package store

import (
	"math"
	"time"
)

func (s *CacheStore) GetFloat32(key string) (float32, error) {
	if v, err := s.GetUInt32(key); err != nil {
		return 0, err
	} else {
		return math.Float32frombits(v), nil
	}
}

func (s *CacheStore) SetFloat32(key string, value float32, exp time.Duration) error {
	bits := math.Float32bits(value)
	return s.SetUInt32(key, bits, exp)
}

// float64
func (s *CacheStore) GetFloat64(key string) (float64, error) {
	if v, err := s.GetUInt64(key); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(v), nil
	}
}

func (s *CacheStore) SetFloat64(key string, value float64, exp time.Duration) error {
	bits := math.Float64bits(value)
	return s.SetUInt64(key, bits, exp)
}
