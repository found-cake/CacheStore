package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetBool(key string) (bool, error) {
	if key == "" {
		return false, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return false, err
	}
	if e.Type != types.BOOLEAN {
		return false, errors.ErrTypeMismatch(key, types.BOOLEAN, e.Type)
	}
	return len(e.Data) > 0 && e.Data[0] == 1, nil
}

func (s *CacheStore) SetBool(key string, value bool, exp time.Duration) error {
	v := byte(0)
	if value {
		v = 1
	}
	return s.Set(key, types.BOOLEAN, []byte{v}, exp)
}
