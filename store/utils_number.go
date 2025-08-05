package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/generic"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) setKeepExp(key string, dataType types.DataType, value []byte, expiry int64) {
	s.memorydbTemporary[key] = entry.Entry{
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
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(expected, e.Type)
	}
	return utils.Binary2UInt16(e.Data)
}

func (s *CacheStore) getNum32(key string, expected types.DataType) (uint32, error) {
	if key == "" {
		return 0, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(expected, e.Type)
	}
	return utils.Binary2UInt32(e.Data)
}

func (s *CacheStore) getNum64(key string, expected types.DataType) (uint64, error) {
	if key == "" {
		return 0, errors.ErrKeyEmpty
	}
	s.temporaryMux.RLock()
	defer s.temporaryMux.RUnlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return 0, err
	}
	if e.Type != expected {
		return 0, errors.ErrTypeMismatch(expected, e.Type)
	}
	return utils.Binary2UInt64(e.Data)
}

func incrNumber[T generic.Numberic](
	s *CacheStore,
	key string,
	delta T,
	data_type types.DataType,
	exp time.Duration,
	fromBinary func([]byte) (T, error),
	toBinary func(T) []byte,
	checkOverFlow func(T, T) bool,
	checkFloatSpesial func(T) bool,
) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.temporaryMux.Lock()
	defer s.temporaryMux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := toBinary(delta)
		s.unsafeSet(key, data_type, data, exp)
		return nil
	}
	if e.Type != data_type {
		return errors.ErrTypeMismatch(data_type, e.Type)
	}
	value, err := fromBinary(e.Data)
	if err != nil {
		return err
	}
	if checkOverFlow(value, delta) {
		return errors.ErrValueOverflow(key, data_type, value, delta)
	}
	value += delta
	data := toBinary(value)
	if checkFloatSpesial != nil && checkFloatSpesial(value) {
		return errors.ErrFloatSpecial
	}
	if exp > 0 {
		s.unsafeSet(key, data_type, data, exp)
	} else {
		s.setKeepExp(key, data_type, data, e.Expiry)
	}
	return nil
}
