package store

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/sqlite"
	"github.com/found-cake/CacheStore/utils/types"
)

type CacheStore struct {
	persistentMux      sync.RWMutex
	temporaryMux       sync.RWMutex
	memorydbPersistent map[string]entry.Entry
	memorydbTemporary  map[string]entry.Entry
	dirty              *dirtyManager
	sqlitedb           *sqlite.SqliteStore
	done               chan struct{}
	wg                 sync.WaitGroup
	closed             atomic.Bool
}

const (
	TTLNoExpiry time.Duration = -1 // Key exists and does not expire
	TTLExpired  time.Duration = -2 // Key does not exist or is expired
)

func (s *CacheStore) cleanExpired() {
	now := time.Now().UnixMilli()

	s.temporaryMux.Lock()
	defer s.temporaryMux.Unlock()

	for key, entry := range s.memorydbTemporary {
		if entry.IsExpiredWithUnixMilli(now) {
			delete(s.memorydbTemporary, key)
		}
	}
}

func (s *CacheStore) unsafeGet(key string) (entry.Entry, error) {
	v, ok := s.memorydbTemporary[key]
	if !ok {
		return v, errors.ErrNoDataForKey(key)
	}
	if v.IsExpired() {
		return v, errors.ErrNoDataForKey(key)
	}
	return v, nil
}

func (s *CacheStore) Get(key string) (types.DataType, []byte, error) {
	if key == "" {
		return types.UNKNOWN, nil, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	v, err := s.unsafeGet(key)
	if err != nil {
		return types.UNKNOWN, nil, err
	}

	result := make([]byte, len(v.Data))
	copy(result, v.Data)
	return v.Type, result, nil
}

// ⚠️  WARNING: GetNoCopy returns a reference to internal cache data.
//
// In contrast to the standard Get() method, which returns a **safe copy**,
//
// GetNoCopy is designed for performance-critical scenarios where copying is avoided.
// However, modifying the returned value may cause unexpected behavior in concurrent environments.
//
// ✅ If you don't explicitly need zero-copy performance,
//
//	use Get() to avoid race conditions and data corruption.
func (s *CacheStore) GetNoCopy(key string) (types.DataType, []byte, error) {
	if key == "" {
		return types.UNKNOWN, nil, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	v, err := s.unsafeGet(key)
	if err != nil {
		return types.UNKNOWN, nil, err
	}

	return v.Type, v.Data, nil
}

func (s *CacheStore) unsafeSet(key string, dataType types.DataType, value []byte, expiry time.Duration) {
	s.memorydbTemporary[key] = entry.NewEntry(dataType, value, expiry)

	if s.dirty != nil {
		s.dirty.set(key)
	}
}

func (s *CacheStore) Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if value == nil {
		return errors.ErrValueNil
	}

	s.temporaryMux.Lock()
	s.memorydbTemporary[key] = entry.NewEntry(dataType, value, expiry)
	s.temporaryMux.Unlock()

	if s.dirty != nil {
		s.dirty.set(key)
	}

	return nil
}

func (s *CacheStore) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	s.temporaryMux.Lock()
	delete(s.memorydbTemporary, key)
	s.temporaryMux.Unlock()

	if s.dirty != nil {
		s.dirty.delete(key)
	}

	return nil
}

func (s *CacheStore) Flush() {
	s.temporaryMux.Lock()
	s.memorydbTemporary = make(map[string]entry.Entry)
	s.temporaryMux.Unlock()
	if s.dirty != nil {
		s.dirty.wantFullSync()
	}
}

func (s *CacheStore) IsClosed() bool {
	return s.closed.Load()
}

func (s *CacheStore) Close() error {
	if s.IsClosed() {
		return nil
	}

	s.closed.Store(true)

	close(s.done)
	s.wg.Wait()

	var err error
	if s.sqlitedb != nil {
		defer func() {
			if err := s.sqlitedb.Close(); err != nil {
				log.Println(err)
			}
		}()
		err = s.sqlitedb.Save(s.memorydbTemporary, true)
	}

	s.memorydbTemporary = nil
	s.dirty = nil

	return err
}

func (s *CacheStore) Exists(keys ...string) int {
	now := time.Now().UnixMilli()
	count := 0

	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()

	for _, key := range keys {
		if e, ok := s.memorydbTemporary[key]; ok {
			if !e.IsExpiredWithUnixMilli(now) {
				count++
			}
		}
	}
	return count
}

func (s *CacheStore) Keys() []string {
	now := time.Now().UnixMilli()
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()

	keys := make([]string, 0, len(s.memorydbTemporary))
	for key, e := range s.memorydbTemporary {
		if !e.IsExpiredWithUnixMilli(now) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s *CacheStore) TTL(key string) time.Duration {
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()

	e, ok := s.memorydbTemporary[key]
	if !ok {
		return TTLExpired
	}

	if e.Expiry == 0 {
		return TTLNoExpiry
	}

	now := time.Now().UnixMilli()
	if now >= e.Expiry {
		return TTLExpired
	}

	remaining := time.Duration(e.Expiry-now) * time.Millisecond
	return remaining
}

func (s *CacheStore) Sync() {
	if s.sqlitedb == nil {
		return
	}
	if s.dirty == nil {
		s.FullSync()
		return
	}

	s.dirty.mux.Lock()
	if s.dirty.needFullSync {
		s.dirty.needFullSync = false
		s.dirty.mux.Unlock()
		s.FullSync()
		return
	}

	dirtySize := s.dirty.size()
	if dirtySize == 0 {
		s.dirty.mux.Unlock()
		return
	}

	s.temporaryMux.RLock()
	if dirtySize > s.dirty.ThresholdCount && dirtySize > int(float64(len(s.memorydbTemporary))*s.dirty.ThresholdRatio) {
		s.temporaryMux.RUnlock()
		s.dirty.mux.Unlock()
		s.FullSync()
		return
	}

	set_keys, delete_keys := s.dirty.keys()
	new_data := make(map[string]entry.Entry, len(set_keys))
	for _, key := range set_keys {
		if e, ok := s.memorydbTemporary[key]; ok {
			dataCopy := make([]byte, len(e.Data))
			copy(dataCopy, e.Data)

			new_data[key] = entry.Entry{
				Type:   e.Type,
				Data:   dataCopy,
				Expiry: e.Expiry,
			}
		}
	}

	s.temporaryMux.RUnlock()
	s.dirty.unsafeClear()
	s.dirty.mux.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.sqlitedb.SaveDirtyData(new_data, delete_keys); err != nil {
			log.Println(err)
		}
	}()
}

func (s *CacheStore) FullSync() {
	if s.sqlitedb == nil {
		return
	}

	s.temporaryMux.RLock()
	snapshot := make(map[string]entry.Entry, len(s.memorydbTemporary))
	for key, e := range s.memorydbTemporary {
		dataCopy := make([]byte, len(e.Data))
		copy(dataCopy, e.Data)

		snapshot[key] = entry.Entry{
			Type:   e.Type,
			Data:   dataCopy,
			Expiry: e.Expiry,
		}
	}
	s.temporaryMux.RUnlock()
	if s.dirty != nil {
		s.dirty.clear()
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.sqlitedb.Save(snapshot, false); err != nil {
			log.Println(err)
		}
	}()
}
