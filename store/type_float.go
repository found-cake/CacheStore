package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetFloat32(key string) (float32, error) {
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data float32, err error) {
		data, err = e.AsFloat32()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
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
	_, data, err := get(s, key, func(e *entry.Entry) (t types.DataType, data float64, err error) {
		data, err = e.AsFloat64()
		if err == nil {
			t = e.Type
		}
		return
	})
	return data, err
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
