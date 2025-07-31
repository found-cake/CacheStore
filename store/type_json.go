package store

import (
	"encoding/json"
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetJSON(key string, target interface{}) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return err
	}
	if e.Type != types.JSON {
		return errors.ErrTypeMismatch(key, types.JSON, e.Type)
	}
	if len(e.Data) == 0 {
		return errors.ErrNoDataForKey(key)
	}
	return json.Unmarshal(e.Data, target)
}

func (s *CacheStore) SetJSON(key string, value interface{}, exp time.Duration) error {
	if data, err := json.Marshal(value); err != nil {
		return err
	} else {
		return s.Set(key, types.JSON, data, exp)
	}
}
