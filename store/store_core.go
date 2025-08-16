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
	memorydbPersistent map[string]entry.Entry
	temporaryMux       sync.RWMutex
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

type entryProcessor[T interface{}] func(*entry.Entry) (types.DataType, T, error)

func get[T interface{}](s *CacheStore, key string, proc entryProcessor[T]) (t types.DataType, data T, err error) {
	if key == "" {
		err = errors.ErrKeyEmpty
		return
	}

	{
		s.persistentMux.RLock()
		defer s.persistentMux.RUnlock()
		v, ok := s.memorydbPersistent[key]
		if ok {
			return proc(&v)
		}
	}

	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	v, ok := s.memorydbTemporary[key]
	if !ok {
		err = errors.ErrNoDataForKey(key)
		return
	}
	if v.IsExpired() {
		err = errors.ErrNoDataForKey(key)
		return
	}
	return proc(&v)
}

func (s *CacheStore) Get(key string) (types.DataType, []byte, error) {
	return get(s, key, func(e *entry.Entry) (types.DataType, []byte, error) {
		return e.Type, e.CopyData(), nil
	})
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
	return get(s, key, func(e *entry.Entry) (types.DataType, []byte, error) {
		return e.Type, e.Data, nil
	})
}

func (s *CacheStore) Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if value == nil {
		return errors.ErrValueNil
	}
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		return tx.Set(key, entry.NewEntry(dataType, value, expiry))
	})
}

func (s *CacheStore) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	s.persistentMux.Lock()
	delete(s.memorydbPersistent, key)
	s.persistentMux.Unlock()

	s.temporaryMux.Lock()
	delete(s.memorydbTemporary, key)
	s.temporaryMux.Unlock()

	if s.dirty != nil {
		s.dirty.delete(key)
	}

	return nil
}

func (s *CacheStore) Flush() {
	s.persistentMux.Lock()
	s.memorydbPersistent = make(map[string]entry.Entry)
	s.persistentMux.Unlock()

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
		s.persistentMux.Lock()
		s.temporaryMux.Lock()
		defer s.persistentMux.Unlock()
		defer s.temporaryMux.Unlock()
		for key, v := range s.memorydbPersistent {
			s.memorydbTemporary[key] = entry.Entry{
				Type: v.Type,
				Data: v.Data,
			}
		}
		err = s.sqlitedb.Save(s.memorydbTemporary, true)
	}

	s.memorydbPersistent = nil
	s.memorydbTemporary = nil
	s.dirty = nil

	return err
}

func (s *CacheStore) Exists(keys ...string) int {
	now := time.Now().UnixMilli()
	count := 0

	s.persistentMux.RLock()
	for _, key := range keys {
		if _, ok := s.memorydbPersistent[key]; ok {
			count++
		}
	}
	s.persistentMux.RUnlock()

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
	s.persistentMux.RLock()
	s.temporaryMux.RLock()

	keys := make([]string, 0, len(s.memorydbPersistent)+len(s.memorydbTemporary))
	for key := range s.memorydbPersistent {
		keys = append(keys, key)
	}
	s.persistentMux.RUnlock()

	now := time.Now().UnixMilli()
	defer s.temporaryMux.RUnlock()

	for key, e := range s.memorydbTemporary {
		if !e.IsExpiredWithUnixMilli(now) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s *CacheStore) TTL(key string) time.Duration {
	s.persistentMux.RLock()
	if _, ok := s.memorydbPersistent[key]; ok {
		return TTLNoExpiry
	}
	s.persistentMux.RUnlock()

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

	s.persistentMux.RLock()
	s.temporaryMux.RLock()
	if dirtySize > s.dirty.ThresholdCount && dirtySize > int(float64(len(s.memorydbPersistent)+len(s.memorydbTemporary))*s.dirty.ThresholdRatio) {
		s.persistentMux.RUnlock()
		s.temporaryMux.RUnlock()
		s.dirty.mux.Unlock()
		s.FullSync()
		return
	}

	set_keys, delete_keys := s.dirty.keys()
	new_data := make(map[string]entry.Entry, len(set_keys))
	for _, key := range set_keys {
		if e, ok := s.memorydbPersistent[key]; ok {
			new_data[key] = entry.Entry{
				Type: e.Type,
				Data: e.CopyData(),
			}
			continue
		}
		if e, ok := s.memorydbTemporary[key]; ok {
			new_data[key] = entry.Entry{
				Type:   e.Type,
				Data:   e.CopyData(),
				Expiry: e.Expiry,
			}
		}
	}

	s.persistentMux.RUnlock()
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

	tx := newSnapshotReadTX(s)
	if s.dirty != nil {
		s.dirty.clear()
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.sqlitedb.Save(tx.memorydb, false); err != nil {
			log.Println(err)
		}
	}()
}
