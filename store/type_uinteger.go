package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetUInt16(key string) (uint16, error) {
	return s.getNum16(key, types.UINT16)
}

func (s *CacheStore) SetUInt16(key string, value uint16, exp time.Duration) error {
	return s.Set(key, types.UINT16, num16tob(value), exp)
}

func (s *CacheStore) IncrUInt16(key string, delta uint16, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := num16tob(delta)
		s.unsafeSet(key, types.UINT16, data, exp)
		return nil
	}
	if e.Type != types.UINT16 {
		return errors.ErrTypeMismatch(key, types.UINT16, e.Type)
	}
	value, err := b2num16(e.Data)
	if err != nil {
		return err
	}
	value += delta
	data := num16tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT16, data, exp)
	} else {
		s.setKeepExp(key, types.UINT16, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) DecrUInt16(key string, delta uint16, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return errors.ErrNoDataForKey(key)
	}
	if e.Type != types.UINT16 {
		return errors.ErrTypeMismatch(key, types.UINT16, e.Type)
	}
	value, err := b2num16(e.Data)
	if err != nil {
		return err
	}
	if value < delta {
		return errors.ErrUnsignedUnderflow(key, value, delta)
	}
	value -= delta
	data := num16tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT16, data, exp)
	} else {
		s.setKeepExp(key, types.UINT16, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) GetUInt32(key string) (uint32, error) {
	return s.getNum32(key, types.UINT32)
}

func (s *CacheStore) SetUInt32(key string, value uint32, exp time.Duration) error {
	return s.Set(key, types.UINT32, num32tob(value), exp)
}

func (s *CacheStore) IncrUInt32(key string, delta uint32, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := num32tob(delta)
		s.unsafeSet(key, types.UINT32, data, exp)
		return nil
	}
	if e.Type != types.UINT32 {
		return errors.ErrTypeMismatch(key, types.UINT32, e.Type)
	}
	value, err := b2num32(e.Data)
	if err != nil {
		return err
	}
	value += delta
	data := num32tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT32, data, exp)
	} else {
		s.setKeepExp(key, types.UINT32, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) DecrUInt32(key string, delta uint32, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return errors.ErrNoDataForKey(key)
	}
	if e.Type != types.UINT32 {
		return errors.ErrTypeMismatch(key, types.UINT32, e.Type)
	}
	value, err := b2num32(e.Data)
	if err != nil {
		return err
	}
	if value < delta {
		return errors.ErrUnsignedUnderflow(key, value, delta)
	}
	value -= delta
	data := num32tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT32, data, exp)
	} else {
		s.setKeepExp(key, types.UINT32, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) GetUInt64(key string) (uint64, error) {
	return s.getNum64(key, types.UINT64)
}

func (s *CacheStore) SetUInt64(key string, value uint64, exp time.Duration) error {
	return s.Set(key, types.UINT64, num64tob(value), exp)
}

func (s *CacheStore) IncrUInt64(key string, delta uint64, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		data := num64tob(delta)
		s.unsafeSet(key, types.UINT64, data, exp)
		return nil
	}
	if e.Type != types.UINT64 {
		return errors.ErrTypeMismatch(key, types.UINT64, e.Type)
	}
	value, err := b2num64(e.Data)
	if err != nil {
		return err
	}
	value += delta
	data := num64tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT64, data, exp)
	} else {
		s.setKeepExp(key, types.UINT64, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) DecrUInt64(key string, delta uint64, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil {
		return errors.ErrNoDataForKey(key)
	}
	if e.Type != types.UINT64 {
		return errors.ErrTypeMismatch(key, types.UINT64, e.Type)
	}
	value, err := b2num64(e.Data)
	if err != nil {
		return err
	}
	if value < delta {
		return errors.ErrUnsignedUnderflow(key, value, delta)
	}
	value -= delta
	data := num64tob(value)
	if exp > 0 {
		s.unsafeSet(key, types.UINT64, data, exp)
	} else {
		s.setKeepExp(key, types.UINT64, data, e.Expiry)
	}
	return nil
}
