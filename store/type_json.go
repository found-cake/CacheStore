package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetJSON(key string, target interface{}) error {
	_, _, err := get(s, key, func(e *entry.Entry) (types.DataType, struct{}, error) {
		err := e.AsJSON(target)
		if err != nil {
			return types.UNKNOWN, struct{}{}, err
		}

		return e.Type, struct{}{}, nil
	})
	return err
}

func (s *CacheStore) SetJSON(key string, value interface{}, exp time.Duration) error {
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		if e, err := entry.FromJSON(value, exp); err == nil {
			return tx.Set(key, e)
		} else {
			return err
		}
	})
}
