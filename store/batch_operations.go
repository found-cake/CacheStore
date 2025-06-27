package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
)

type BatchItem struct {
	Key    string
	Value  []byte
	Expiry time.Duration
}

type BatchResult struct {
	Key   string
	Value []byte
	Error error
}

func (s *CacheStore) MGet(keys ...string) []BatchResult {
	if len(keys) == 0 {
		return nil
	}

	results := make([]BatchResult, len(keys))
	now := uint32(time.Now().Unix())

	s.mux.RLock()
	defer s.mux.RUnlock()

	for i, key := range keys {
		results[i].Key = key
		if key == "" {
			results[i].Error = errors.ErrKeyEmpty
			continue
		}
		if e, ok := s.memorydb[key]; ok {
			if !e.IsExpiredWithTime(now) {
				results[i].Value = e.Data
			} else {
				results[i].Error = errors.ErrNoDataForKey(key)
			}
		} else {
			results[i].Error = errors.ErrNoDataForKey(key)
		}
	}

	return results
}

func (s *CacheStore) MSet(items ...BatchItem) []error {
	if len(items) == 0 {
		return nil
	}

	errs := make([]error, len(items))

	s.mux.Lock()
	defer s.mux.Unlock()

	for i, item := range items {
		if item.Key == "" {
			errs[i] = errors.ErrKeyEmpty
			continue
		}
		if item.Value == nil {
			errs[i] = errors.ErrValueNil
			continue
		}

		s.memorydb[item.Key] = entry.NewEntry(item.Value, item.Expiry)
	}

	return errs
}

func (s *CacheStore) MDelete(keys ...string) []error {
	if len(keys) == 0 {
		return nil
	}

	errs := make([]error, len(keys))

	s.mux.Lock()
	defer s.mux.Unlock()

	for i, key := range keys {
		if key == "" {
			errs[i] = errors.ErrKeyEmpty
			continue
		}

		delete(s.memorydb, key)
	}

	return errs
}
