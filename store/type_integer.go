package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetInt16(key string) (int16, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data int16, err error) {
		data, err = e.AsInt16()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetInt16(key string, value int16, exp time.Duration) error {
	return s.Set(key, types.INT16, utils.Int16toBinary(value), exp)
}

func (s *CacheStore) IncrInt16(key string, delta int16, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.INT16, exp,
		utils.Binary2Int16,
		utils.Int16toBinary,
		utils.Int16CheckOver,
		nil,
	)
}

func (s *CacheStore) GetInt32(key string) (int32, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data int32, err error) {
		data, err = e.AsInt32()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetInt32(key string, value int32, exp time.Duration) error {
	return s.Set(key, types.INT32, utils.Int32toBinary(value), exp)
}

func (s *CacheStore) IncrInt32(key string, delta int32, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.INT32, exp,
		utils.Binary2Int32,
		utils.Int32toBinary,
		utils.Int32CheckOver,
		nil,
	)
}

func (s *CacheStore) GetInt64(key string) (int64, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data int64, err error) {
		data, err = e.AsInt64()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
}

func (s *CacheStore) SetInt64(key string, value int64, exp time.Duration) error {
	return s.Set(key, types.INT64, utils.Int64toBinary(value), exp)
}

func (s *CacheStore) IncrInt64(key string, delta int64, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.INT64, exp,
		utils.Binary2Int64,
		utils.Int64toBinary,
		utils.Int64CheckOver,
		nil,
	)
}
