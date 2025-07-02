package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
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

func (s *CacheStore) IncrInt16(key string, delta int16, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil || e.IsExpired() {
		data := num16tob(uint16(delta))
		s.unsafeSet(key, types.INT16, data, exp)
		return nil
	}
	if e.Type != types.INT16 {
		return errors.ErrTypeMismatch(key, types.INT16, e.Type)
	}
	uvalue, err := b2num16(e.Data)
	if err != nil {
		return err
	}
	value := int16(uvalue) + delta
	data := num16tob(uint16(value))
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
	return s.Set(key, types.INT32, num32tob(uint32(value)), exp)
}

func (s *CacheStore) IncrInt32(key string, delta int32, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil || e.IsExpired() {
		data := num32tob(uint32(delta))
		s.unsafeSet(key, types.INT32, data, exp)
		return nil
	}
	if e.Type != types.INT32 {
		return errors.ErrTypeMismatch(key, types.INT32, e.Type)
	}
	uvalue, err := b2num32(e.Data)
	if err != nil {
		return err
	}
	value := int32(uvalue) + delta
	data := num32tob(uint32(value))
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
	return s.Set(key, types.INT64, num64tob(uint64(value)), exp)
}

func (s *CacheStore) IncrInt64(key string, delta int64, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil || e.IsExpired() {
		data := num64tob(uint64(delta))
		s.unsafeSet(key, types.INT16, data, exp)
		return nil
	}
	if e.Type != types.INT64 {
		return errors.ErrTypeMismatch(key, types.INT64, e.Type)
	}
	uvalue, err := b2num64(e.Data)
	if err != nil {
		return err
	}
	value := int64(uvalue) + delta
	data := num64tob(uint64(value))
	if exp > 0 {
		s.unsafeSet(key, types.INT64, data, exp)
	} else {
		s.setKeepExp(key, types.INT64, data, e.Expiry)
	}
	return nil
}
