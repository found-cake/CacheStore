package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetString(key string) (string, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data string, err error) {
		data, err = e.AsString()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetString(key string, value string, exp time.Duration) error {
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		return tx.Set(key, entry.FromString(value, exp))
	})
}
