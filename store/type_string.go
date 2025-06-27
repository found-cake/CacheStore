package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetString(key string) (string, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return "", err
	}
	if t != types.STRING {
		return "", errors.ErrTypeMismatch(key, types.STRING, t)
	}
	return string(data), nil
}

func (s *CacheStore) SetString(key string, value string, exp time.Duration) error {
	return s.Set(key, types.STRING, []byte(value), exp)
}
