package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetTime(key string) (time.Time, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data time.Time, err error) {
		data, err = e.AsTime()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetTime(key string, value time.Time, exp time.Duration) error {
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		if e, err := entry.FromTime(value, exp); err == nil {
			return tx.Set(key, e)
		} else {
			return err
		}
	})
}
