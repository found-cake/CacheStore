package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetRaw(key string) ([]byte, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data []byte, err error) {
		data, err = e.AsRaw()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) GetRawNoCopy(key string) ([]byte, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data []byte, err error) {
		data, err = e.AsRawNoCopy()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetRaw(key string, value []byte, exp time.Duration) error {
	return s.WriteTransaction(func(tx *WriteTransaction) error {
		return tx.Set(key, entry.FromRaw(value, exp))
	})
}
