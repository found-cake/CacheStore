package store

import (
	"math"
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
)

func (s *CacheStore) GetFloat32(key string) (float32, error) {
	if v, err := s.getNum32(key, types.FLOAT32); err != nil {
		return 0, err
	} else {
		return math.Float32frombits(v), nil
	}
}

func (s *CacheStore) SetFloat32(key string, value float32, exp time.Duration) error {
	bits := math.Float32bits(value)
	return s.Set(key, types.FLOAT32, num32tob(bits), exp)
}

func (s *CacheStore) IncrFloat32(key string, delta float32, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil || e.IsExpired() {
		data := num32tob(math.Float32bits(delta))
		s.unsafeSet(key, types.FLOAT32, data, exp)
		return nil
	}
	if e.Type != types.FLOAT32 {
		return errors.ErrTypeMismatch(key, types.FLOAT32, e.Type)
	}
	uvalue, err := b2num32(e.Data)
	if err != nil {
		return err
	}
	value := math.Float32frombits(uvalue) + delta
	data := num32tob(math.Float32bits(value))
	if exp > 0 {
		s.unsafeSet(key, types.FLOAT32, data, exp)
	} else {
		s.setKeepExp(key, types.FLOAT32, data, e.Expiry)
	}
	return nil
}

func (s *CacheStore) GetFloat64(key string) (float64, error) {
	if v, err := s.getNum64(key, types.FLOAT64); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(v), nil
	}
}

func (s *CacheStore) SetFloat64(key string, value float64, exp time.Duration) error {
	bits := math.Float64bits(value)
	return s.Set(key, types.FLOAT64, num64tob(bits), exp)
}

func (s *CacheStore) IncrFloat64(key string, delta float64, exp time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	e, err := s.unsafeGet(key)
	if err != nil || e.IsExpired() {
		data := num64tob(math.Float64bits(delta))
		s.unsafeSet(key, types.FLOAT64, data, exp)
		return nil
	}
	if e.Type != types.FLOAT64 {
		return errors.ErrTypeMismatch(key, types.FLOAT64, e.Type)
	}
	uvalue, err := b2num64(e.Data)
	if err != nil {
		return err
	}
	value := math.Float64frombits(uvalue) + delta
	data := num64tob(math.Float64bits(value))
	if exp > 0 {
		s.unsafeSet(key, types.FLOAT64, data, exp)
	} else {
		s.setKeepExp(key, types.FLOAT64, data, e.Expiry)
	}
	return nil
}
