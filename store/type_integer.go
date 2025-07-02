package store

import (
	"time"

	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetInt16(key string) (int16, error) {
	if v, err := s.getNum16(key, types.INT16); err != nil {
		return 0, err
	} else {
		return int16(v), nil
	}
}

func (s *CacheStore) SetInt16(key string, value int16, exp time.Duration) error {
	return s.Set(key, types.INT16, num16tob(uint16(value)), exp)
}

func (s *CacheStore) GetInt32(key string) (int32, error) {
	if v, err := s.getNum32(key, types.INT32); err != nil {
		return 0, err
	} else {
		return int32(v), nil
	}
}

func (s *CacheStore) SetInt32(key string, value int32, exp time.Duration) error {
	return s.Set(key, types.INT32, num32tob(uint32(value)), exp)
}

func (s *CacheStore) GetInt64(key string) (int64, error) {
	if v, err := s.getNum64(key, types.INT64); err != nil {
		return 0, err
	} else {
		return int64(v), nil
	}
}

func (s *CacheStore) SetInt64(key string, value int64, exp time.Duration) error {
	return s.Set(key, types.INT64, num64tob(uint64(value)), exp)
}
