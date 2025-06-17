package cachestore

import (
	"sync"
	"time"
)

type entry struct {
	data   []byte
	expiry uint32
}

type CacheStore struct {
	mux      sync.RWMutex
	memorydb map[string]entry
	done     chan bool
	config   Config
}

const (
	TTLNoExpiry time.Duration = -1 // Key exists and does not expire
	TTLExpired  time.Duration = -2 // Key does not exist or is expired
)

func NewCacheStore(cfg Config) (*CacheStore, error) {
	store := &CacheStore{
		memorydb: make(map[string]entry),
		done:     make(chan bool),
		config:   cfg,
	}
	if cfg.DBSave {
		if cfg.DBFileName == "" {
			cfg.DBFileName = "cache.db"
		}
		db, err := initDB(cfg.DBFileName)
		if err != nil {
			return nil, err
		}
		defer db.Close()
		if data, err := loadFromDB(db); err != nil {
			return nil, err
		} else {
			store.memorydb = data
		}
	}

	if cfg.GCInterval > 0 {
		go store.gc()
	}

	return store, nil
}

func (s *CacheStore) gc() {
	ticker := time.NewTicker(s.config.GCInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanExpired()
		case <-s.done:
			return
		}
	}
}

func (s *CacheStore) cleanExpired() {
	now := uint32(time.Now().Unix())
	s.mux.Lock()
	defer s.mux.Unlock()

	for key, entry := range s.memorydb {
		if entry.expiry > 0 && entry.expiry <= now {
			delete(s.memorydb, key)
		}
	}
}

func (s *CacheStore) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, ErrKeyEmpty
	}
	s.mux.RLock()
	v, ok := s.memorydb[key]
	s.mux.RUnlock()
	if !ok {
		return nil, nil
	}
	if v.expiry > 0 && v.expiry <= uint32(time.Now().Unix()) {
		return nil, nil
	}
	return v.data, nil
}

func (s *CacheStore) MGet(keys ...string) [][]byte {
	result := make([][]byte, len(keys))
	now := uint32(time.Now().Unix())

	s.mux.RLock()
	defer s.mux.RUnlock()

	for i, key := range keys {
		if entry, ok := s.memorydb[key]; ok {
			if entry.expiry == 0 || entry.expiry > now {
				result[i] = entry.data
			}
		}
	}

	return result
}

func (s *CacheStore) Set(key string, value []byte, exp time.Duration) error {
	if key == "" {
		return ErrKeyEmpty
	}
	if value == nil {
		return ErrValueNil
	}

	var expiry uint32
	if exp > 0 {
		expiry = uint32(time.Now().Add(exp).Unix())
	}

	s.mux.Lock()
	s.memorydb[key] = entry{
		data:   value,
		expiry: expiry,
	}
	s.mux.Unlock()

	return nil
}

func (s *CacheStore) Delete(key string) error {
	if key == "" {
		return ErrKeyEmpty
	}

	s.mux.Lock()
	delete(s.memorydb, key)
	s.mux.Unlock()

	return nil
}

func (s *CacheStore) Flush() {
	s.mux.Lock()
	s.memorydb = make(map[string]entry)
	s.mux.Unlock()
}

func (s *CacheStore) Close() error {
	select {
	case <-s.done:
		return nil
	default:
		close(s.done)
	}

	if s.config.DBSave {
		db, err := initDB(s.config.DBFileName)
		if err != nil {
			return err
		}
		defer db.Close()
		return saveDB(db, s.memorydb)
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
			if e.expiry <= 0 || e.expiry > now {
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
		if e.expiry == 0 || e.expiry > now {
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

	if e.expiry <= 0 {
		return TTLNoExpiry
	}

	now := uint32(time.Now().Unix())
	if now >= e.expiry {
		return TTLExpired
	}

	remaining := time.Duration(e.expiry-now) * time.Second
	return remaining
}
