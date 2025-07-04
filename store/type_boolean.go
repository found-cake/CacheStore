package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetBool(key string) (bool, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return false, err
	}
	if t != types.BOOLEAN {
		return false, errors.ErrTypeMismatch(key, types.BOOLEAN, t)
	}
	return len(data) > 0 && data[0] == 1, nil
}

func (s *CacheStore) SetBool(key string, value bool, exp time.Duration) error {
	v := byte(0)
	if value {
		v = 1
	}
	return s.Set(key, types.BOOLEAN, []byte{v}, exp)
}
