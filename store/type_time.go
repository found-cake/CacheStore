package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetTime(key string) (time.Time, error) {
	var t time.Time
	if key == "" {
		return t, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return t, err
	}
	if e.Type != types.TIME {
		return t, errors.ErrTypeMismatch(key, types.TIME, e.Type)
	}
	if len(e.Data) == 0 {
		return t, errors.ErrNoDataForKey(key)
	}
	err = t.UnmarshalBinary(e.Data)
	return t, err
}

func (s *CacheStore) SetTime(key string, value time.Time, exp time.Duration) error {
	if b, err := value.MarshalBinary(); err != nil {
		return err
	} else {
		return s.Set(key, types.TIME, b, exp)
	}
}
