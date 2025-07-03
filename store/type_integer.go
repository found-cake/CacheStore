package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
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
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := utils.Int16toBinary(delta)
		s.unsafeSet(key, types.INT16, data, exp)
		return nil
	}
	if e.Type != types.INT16 {
		return errors.ErrTypeMismatch(key, types.INT16, e.Type)
	}
	value, err := utils.Binary2Int16(e.Data)
	if err != nil {
		return err
	}
	if utils.Int16CheckOver(value, delta) {
		return errors.ErrValueOverflow(key, types.INT16, value, delta)
	}
	value += delta
	data := utils.Int16toBinary(value)
	if exp > 0 {
		s.unsafeSet(key, types.INT16, data, exp)
	} else {
		s.setKeepExp(key, types.INT16, data, e.Expiry)
	}
	return nil
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
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := utils.Int32toBinary(delta)
		s.unsafeSet(key, types.INT32, data, exp)
		return nil
	}
	if e.Type != types.INT32 {
		return errors.ErrTypeMismatch(key, types.INT32, e.Type)
	}
	value, err := utils.Binary2Int32(e.Data)
	if err != nil {
		return err
	}
	if utils.Int32CheckOver(value, delta) {
		return errors.ErrValueOverflow(key, types.INT32, value, delta)
	}
	value += delta
	data := utils.Int32toBinary(value)
	if exp > 0 {
		s.unsafeSet(key, types.INT32, data, exp)
	} else {
		s.setKeepExp(key, types.INT32, data, e.Expiry)
	}
	return nil
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
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := utils.Int64toBinary(delta)
		s.unsafeSet(key, types.INT64, data, exp)
		return nil
	}
	if e.Type != types.INT64 {
		return errors.ErrTypeMismatch(key, types.INT64, e.Type)
	}
	value, err := utils.Binary2Int64(e.Data)
	if err != nil {
		return err
	}
	if utils.Int64CheckOver(value, delta) {
		return errors.ErrValueOverflow(key, types.INT64, value, delta)
	}
	value += delta
	data := utils.Int64toBinary(value)
	if exp > 0 {
		s.unsafeSet(key, types.INT64, data, exp)
	} else {
		s.setKeepExp(key, types.INT64, data, e.Expiry)
	}
	return nil
}
