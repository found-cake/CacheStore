package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetRaw(key string) ([]byte, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	if t != types.RAW {
		return nil, errors.ErrTypeMismatch(key, types.RAW, t)
	}
	return data, nil
}

func (s *CacheStore) SetRaw(key string, value []byte, exp time.Duration) error {
	return s.Set(key, types.RAW, value, exp)
}
