package store

import (
	"math"
	"time"

	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetFloat32(key string) (float32, error) {
	if v, err := s.getNum32(key, types.FLOAT32); err != nil {
		return 0, err
	} else {
		return math.Float32frombits(v), nil
	}
}

func (s *CacheStore) SetFloat32(key string, value float32, exp time.Duration) error {
	return s.Set(key, types.FLOAT32, utils.Float32toBinary(value), exp)
}

func (s *CacheStore) IncrFloat32(key string, delta float32, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.FLOAT32, exp,
		utils.Binary2Float32,
		utils.Float32toBinary,
		utils.Float32CheckOver,
		utils.CheckFloat32Special,
	)
}

func (s *CacheStore) GetFloat64(key string) (float64, error) {
	if v, err := s.getNum64(key, types.FLOAT64); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(v), nil
	}
}

func (s *CacheStore) SetFloat64(key string, value float64, exp time.Duration) error {
	return s.Set(key, types.FLOAT64, utils.Float64toBinary(value), exp)
}

func (s *CacheStore) IncrFloat64(key string, delta float64, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.FLOAT64, exp,
		utils.Binary2Float64,
		utils.Float64toBinary,
		utils.Float64CheckOver,
		utils.CheckFloat64Special,
	)
}
