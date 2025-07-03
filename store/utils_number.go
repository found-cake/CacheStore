package store

import (
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
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

func (s *CacheStore) getNum16(key string, expected types.DataType) (uint16, error) {
	if key == "" {
		return 0, errors.ErrKeyEmpty
	}
	s.mux.RLock()
	defer s.mux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(key, expected, e.Type)
	}
	return utils.Binary2UInt16(e.Data)
}

func (s *CacheStore) getNum32(key string, expected types.DataType) (uint32, error) {
	if key == "" {
		return 0, errors.ErrKeyEmpty
	}
	s.mux.RLock()
	defer s.mux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(key, expected, e.Type)
	}
	return utils.Binary2UInt32(e.Data)
}

func (s *CacheStore) getNum64(key string, expected types.DataType) (uint64, error) {
	if key == "" {
		return 0, errors.ErrKeyEmpty
	}
	s.mux.RLock()
	defer s.mux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(key, expected, e.Type)
	}
	return utils.Binary2UInt64(e.Data)
}
