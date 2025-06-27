package store

import (
	"encoding/binary"
	"time"

	"github.com/found-cake/CacheStore/errors"
)

func (s *CacheStore) GetUInt16(key string) (uint16, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 2 {
		return 0, errors.ErrInvalidDataLength(2, len(data))
	}
	return binary.LittleEndian.Uint16(data), nil
}

func (s *CacheStore) SetUInt16(key string, value uint16, exp time.Duration) error {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return s.Set(key, buffer, exp)
}

func (s *CacheStore) GetUInt32(key string) (uint32, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, errors.ErrInvalidDataLength(4, len(data))
	}
	return binary.LittleEndian.Uint32(data), nil
}

func (s *CacheStore) SetUInt32(key string, value uint32, exp time.Duration) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return s.Set(key, buffer, exp)
}

func (s *CacheStore) GetUInt64(key string) (uint64, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 8 {
		return 0, errors.ErrInvalidDataLength(8, len(data))
	}
	return binary.LittleEndian.Uint64(data), nil
}

func (s *CacheStore) SetUInt64(key string, value uint64, exp time.Duration) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return s.Set(key, buffer, exp)
}
