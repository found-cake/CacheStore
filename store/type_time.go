package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetTime(key string) (time.Time, error) {
	var t time.Time
	dt, data, err := s.Get(key)
	if err != nil {
		return t, err
	}
	if dt != types.TIME {
		return t, errors.ErrTypeMismatch(key, types.TIME, dt)
	}
	if len(data) == 0 {
		return t, errors.ErrNoDataForKey(key)
	}
	err = t.UnmarshalBinary(data)
	return t, err
}

func (s *CacheStore) SetTime(key string, value time.Time, exp time.Duration) error {
	if b, err := value.MarshalBinary(); err != nil {
		return err
	} else {
		return s.Set(key, types.TIME, b, exp)
	}
}
