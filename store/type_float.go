package store

import (
	"math"
	"time"

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

func (s *CacheStore) GetFloat64(key string) (float64, error) {
	if v, err := s.getNum64(key, types.FLOAT64); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(v), nil
	}
}

func (s *CacheStore) SetFloat64(key string, value float64, exp time.Duration) error {
	bits := math.Float64bits(value)
	return s.Set(key, types.FLOAT32, num64tob(bits), exp)
}
