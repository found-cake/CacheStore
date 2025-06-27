package store

import (
	"encoding/json"
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetJSON(key string, target interface{}) error {
	t, data, err := s.Get(key)
	if err != nil {
		return err
	}
	if t != types.JSON {
		return errors.ErrTypeMismatch(key, types.JSON, t)
	}
	if len(data) == 0 {
		return errors.ErrNoDataForKey(key)
	}
	return json.Unmarshal(data, target)
}

func (s *CacheStore) SetJSON(key string, value interface{}, exp time.Duration) error {
	if data, err := json.Marshal(value); err != nil {
		return err
	} else {
		return s.Set(key, types.JSON, data, exp)
	}
}
