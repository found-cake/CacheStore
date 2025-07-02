package store

import (
	"encoding/binary"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) setKeepExp(key string, dataType types.DataType, value []byte, expiry uint32) {
	s.memorydb[key] = entry.Entry{
		Type:   dataType,
		Data:   value,
		Expiry: expiry,
	}
	if s.dirty != nil {
		s.dirty.set(key)
	}
}

func b2num16(data []byte) (uint16, error) {
	if len(data) != 2 {
		return 0, errors.ErrInvalidDataLength(2, len(data))
	}
	return binary.LittleEndian.Uint16(data), nil
}

func num16tob(value uint16) []byte {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return buffer
}

func (s *CacheStore) getNum16(key string, expected types.DataType) (uint16, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	return b2num16(data)
}

func b2num32(data []byte) (uint32, error) {
	if len(data) != 4 {
		return 0, errors.ErrInvalidDataLength(4, len(data))
	}
	return binary.LittleEndian.Uint32(data), nil
}

func num32tob(value uint32) []byte {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return buffer
}

func (s *CacheStore) getNum32(key string, expected types.DataType) (uint32, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	return b2num32(data)
}

func b2num64(data []byte) (uint64, error) {
	if len(data) != 8 {
		return 0, errors.ErrInvalidDataLength(8, len(data))
	}
	return binary.LittleEndian.Uint64(data), nil
}

func num64tob(value uint64) []byte {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return buffer
}

func (s *CacheStore) getNum64(key string, expected types.DataType) (uint64, error) {
	t, data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if t != expected {
		return 0, errors.ErrTypeMismatch(key, expected, t)
	}
	return b2num64(data)
}
