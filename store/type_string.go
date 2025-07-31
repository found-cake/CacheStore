package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetString(key string) (string, error) {
	if key == "" {
		return "", errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return "", err
	}
	if e.Type != types.STRING {
		return "", errors.ErrTypeMismatch(key, types.STRING, e.Type)
	}
	return string(e.Data), nil
}

func (s *CacheStore) SetString(key string, value string, exp time.Duration) error {
	return s.Set(key, types.STRING, []byte(value), exp)
}
