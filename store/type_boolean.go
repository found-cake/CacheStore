package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetBool(key string) (bool, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data bool, err error) {
		data, err = e.AsBool()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetBool(key string, value bool, exp time.Duration) error {
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		return tx.Set(key, entry.FromBool(value, exp))
	})
}
