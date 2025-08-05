package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetRaw(key string) ([]byte, error) {
	if key == "" {
		return nil, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return nil, err
	}
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	result := make([]byte, len(e.Data))
	copy(result, e.Data)

	return result, nil
}

func (s *CacheStore) GetRawNoCopy(key string) ([]byte, error) {
	if key == "" {
		return nil, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return nil, err
	}
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	return e.Data, nil
}

func (s *CacheStore) SetRaw(key string, value []byte, exp time.Duration) error {
	return s.Set(key, types.RAW, value, exp)
}
