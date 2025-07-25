package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type BatchItem struct {
	Key   string
	Entry *entry.Entry
}

func NewItem(key string, dataType types.DataType, data []byte, expiry time.Duration) BatchItem {
	if data == nil {
		return BatchItem{Key: key}
	}
	entry := entry.NewEntry(dataType, data, expiry)
	return BatchItem{
		Key:   key,
		Entry: &entry,
	}
}

type BatchResult struct {
	Key   string
	Type  types.DataType
	Value []byte
	Error error
}

func (s *CacheStore) MGet(keys ...string) []BatchResult {
	if len(keys) == 0 {
		return nil
	}

	results := make([]BatchResult, len(keys))
	now := time.Now().UnixMilli()

	s.mux.RLock()
	defer s.mux.RUnlock()

	for i, key := range keys {
		results[i].Key = key
		if key == "" {
			results[i].Error = errors.ErrKeyEmpty
			continue
		}
		if e, ok := s.memorydb[key]; ok {
			if !e.IsExpiredWithUnixMilli(now) {
				cData := make([]byte, len(e.Data))
				copy(cData, e.Data)
				results[i].Type = e.Type
				results[i].Value = cData
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
	if s.dirty != nil {
		s.dirty.mux.Lock()
		defer s.dirty.mux.Unlock()
	}

	for i, item := range items {
		if item.Key == "" {
			errs[i] = errors.ErrKeyEmpty
			continue
		}
		if item.Entry == nil {
			errs[i] = errors.ErrValueNil
			continue
		}
		s.memorydb[item.Key] = *item.Entry
		if s.dirty != nil {
			s.dirty.unsafeSet(item.Key)
		}
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
	if s.dirty != nil {
		s.dirty.mux.Lock()
		defer s.dirty.mux.Unlock()
	}

	for i, key := range keys {
		if key == "" {
			errs[i] = errors.ErrKeyEmpty
			continue
		}

		delete(s.memorydb, key)
		if s.dirty != nil {
			s.dirty.unsafeDelete(key)
		}
	}

	return errs
}
