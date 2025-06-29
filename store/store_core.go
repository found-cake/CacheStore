package store

import (
	"log"
	"sync"
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/sqlite"
	"github.com/found-cake/CacheStore/store/types"
)

type CacheStore struct {
	mux      sync.RWMutex
	memorydb map[string]entry.Entry
	sqlitedb *sqlite.SqliteStore
	done     chan struct{}
}

const (
	TTLNoExpiry time.Duration = -1 // Key exists and does not expire
	TTLExpired  time.Duration = -2 // Key does not exist or is expired
)

func (s *CacheStore) createTicker(gcInterval time.Duration, saveInterval time.Duration) func() {
	if gcInterval > 0 && saveInterval > 0 {
		return func() {
			gcticker := time.NewTicker(gcInterval)
			dbticker := time.NewTicker(saveInterval)
			defer gcticker.Stop()
			defer dbticker.Stop()

			for {
				select {
				case <-gcticker.C:
					s.cleanExpired()
				case <-dbticker.C:
					s.Sync()
				case <-s.done:
					return
				}
			}
		}
	}
	if  gcInterval > 0 {
		return func() {
			gcticker := time.NewTicker(gcInterval)
			defer gcticker.Stop()

			for {
				select {
				case <-gcticker.C:
					s.cleanExpired()
				case <-s.done:
					return
				}
			}
		}
	}
	if saveInterval > 0 {
		return func() {
			dbticker := time.NewTicker(saveInterval)
			defer dbticker.Stop()

			for {
				select {
				case <-dbticker.C:
					s.Sync()
				case <-s.done:
					return
				}
			}
		}
	}
	return nil
}

func (s *CacheStore) cleanExpired() {
	now := uint32(time.Now().Unix())

	s.mux.Lock()
	defer s.mux.Unlock()

	for key, entry := range s.memorydb {
		if entry.IsExpiredWithTime(now) {
			delete(s.memorydb, key)
		}
	}
}

func (s *CacheStore) Get(key string) (types.DataType, []byte, error) {
	if key == "" {
		return types.UNKNOWN, nil, errors.ErrKeyEmpty
	}
	s.mux.RLock()
	defer s.mux.RUnlock()
	v, ok := s.memorydb[key]

	if !ok {
		return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
	}
	if v.IsExpired() {
		return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
	}

	result := make([]byte, len(v.Data))
	copy(result, v.Data)
	return v.Type, result, nil
}

func (s *CacheStore) Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if value == nil {
		return errors.ErrValueNil
	}

	s.mux.Lock()
	s.memorydb[key] = entry.NewEntry(dataType, value, expiry)
	s.mux.Unlock()

	return nil
}

func (s *CacheStore) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	s.mux.Lock()
	delete(s.memorydb, key)
	s.mux.Unlock()

	return nil
}

func (s *CacheStore) Flush() {
	s.mux.Lock()
	s.memorydb = make(map[string]entry.Entry)
	s.mux.Unlock()
}

func (s *CacheStore) Close() error {
	select {
	case <-s.done:
		return nil
	default:
		close(s.done)
	}

	if s.sqlitedb != nil {
		defer func() {
			if err := s.sqlitedb.Close(); err != nil {
				log.Panicln(err)
			}
		}()
		return s.sqlitedb.Save(s.memorydb, true)
	}
	return nil
}

func (s *CacheStore) Exists(keys ...string) int {
	now := uint32(time.Now().Unix())
	count := 0

	s.mux.RLock()
	defer s.mux.RUnlock()

	for _, key := range keys {
		if e, ok := s.memorydb[key]; ok {
			if !e.IsExpiredWithTime(now) {
				count++
			}
		}
	}
	return count
}

func (s *CacheStore) Keys() []string {
	now := uint32(time.Now().Unix())
	s.mux.RLock()
	defer s.mux.RUnlock()

	keys := make([]string, 0, len(s.memorydb))
	for key, e := range s.memorydb {
		if !e.IsExpiredWithTime(now) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s *CacheStore) TTL(key string) time.Duration {
	s.mux.RLock()
	defer s.mux.RUnlock()

	e, ok := s.memorydb[key]
	if !ok {
		return TTLExpired
	}

	if e.Expiry == 0 {
		return TTLNoExpiry
	}

	now := uint32(time.Now().Unix())
	if now >= e.Expiry {
		return TTLExpired
	}

	remaining := time.Duration(e.Expiry-now) * time.Second
	return remaining
}
