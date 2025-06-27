package store

import (
	"time"

	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetUInt16(key string) (uint16, error) {
	return s.getNum16(key, types.UINT16)
}

func (s *CacheStore) SetUInt16(key string, value uint16, exp time.Duration) error {
	return s.setNum16(key, types.UINT16, value, exp)
}

func (s *CacheStore) GetUInt32(key string) (uint32, error) {
	return s.getNum32(key, types.UINT32)
}

func (s *CacheStore) SetUInt32(key string, value uint32, exp time.Duration) error {
	return s.setNum32(key, types.UINT32, value, exp)
}

func (s *CacheStore) GetUInt64(key string) (uint64, error) {
	return s.getNum64(key, types.UINT64)
}

func (s *CacheStore) SetUInt64(key string, value uint64, exp time.Duration) error {
	return s.setNum64(key, types.UINT64, value, exp)
}
