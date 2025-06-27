package store

import (
	"encoding/binary"
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) getNum16(key string, expected types.DataType) (uint16, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	if len(data) != 2 {
		return 0, errors.ErrInvalidDataLength(2, len(data))
	}
	return binary.LittleEndian.Uint16(data), nil
}

func (s *CacheStore) setNum16(key string, dataType types.DataType, value uint16, exp time.Duration) error {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return s.Set(key, dataType, buffer, exp)
}

func (s *CacheStore) getNum32(key string, expected types.DataType) (uint32, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	if len(data) != 4 {
		return 0, errors.ErrInvalidDataLength(4, len(data))
	}
	return binary.LittleEndian.Uint32(data), nil
}

func (s *CacheStore) setNum32(key string, dataType types.DataType, value uint32, exp time.Duration) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return s.Set(key, dataType, buffer, exp)
}

func (s *CacheStore) getNum64(key string, expected types.DataType) (uint64, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	if len(data) != 8 {
		return 0, errors.ErrInvalidDataLength(8, len(data))
	}
	return binary.LittleEndian.Uint64(data), nil
}

func (s *CacheStore) setNum64(key string, dataType types.DataType, value uint64, exp time.Duration) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return s.Set(key, dataType, buffer, exp)
}
