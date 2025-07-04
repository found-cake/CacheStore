package store

import (
	"time"

	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetInt16(key string) (int16, error) {
	if v, err := s.getNum16(key, types.INT16); err != nil {
		return 0, err
	} else {
		return int16(v), nil
	}
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
	if v, err := s.getNum32(key, types.INT32); err != nil {
		return 0, err
	} else {
		return int32(v), nil
	}
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
	if v, err := s.getNum64(key, types.INT64); err != nil {
		return 0, err
	} else {
		return int64(v), nil
	}
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
