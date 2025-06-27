package store

import (
	"encoding/json"
	"time"

	"github.com/found-cake/CacheStore/errors"
)

func (s *CacheStore) GetJSON(key string, target interface{}) error {
	data, err := s.Get(key)
	if err != nil {
		return err
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
		return s.Set(key, data, exp)
	}
}
