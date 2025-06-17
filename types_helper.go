package cachestore

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"time"
)

// string
func (s *CacheStore) GetString(key string) (string, error) {
	if data, err := s.Get(key); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (s *CacheStore) SetString(key string, value string, exp time.Duration) error {
	return s.Set(key, []byte(value), exp)
}

// boolean
func (s *CacheStore) GetBool(key string) (bool, error) {
	data, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return len(data) > 0 && data[0] == 1, nil
}

func (s *CacheStore) SetBool(key string, value bool, exp time.Duration) error {
	v := byte(0)
	if value {
		v = 1
	}
	return s.Set(key, []byte{v}, exp)
}

// uint16
func (s *CacheStore) GetUInt16(key string) (uint16, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 2 {
		return 0, ErrInvalidDataLength(2, len(data))
	}
	return binary.LittleEndian.Uint16(data), nil
}

func (s *CacheStore) SetUInt16(key string, value uint16, exp time.Duration) error {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return s.Set(key, buffer, exp)
}

// uint32
func (s *CacheStore) GetUInt32(key string) (uint32, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 4 {
		return 0, ErrInvalidDataLength(4, len(data))
	}
	return binary.LittleEndian.Uint32(data), nil
}

func (s *CacheStore) SetUInt32(key string, value uint32, exp time.Duration) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return s.Set(key, buffer, exp)
}

// uint64
func (s *CacheStore) GetUInt64(key string) (uint64, error) {
	data, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(data) < 8 {
		return 0, ErrInvalidDataLength(8, len(data))
	}
	return binary.LittleEndian.Uint64(data), nil
}

func (s *CacheStore) SetUInt64(key string, value uint64, exp time.Duration) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return s.Set(key, buffer, exp)
}

// int16
func (s *CacheStore) GetInt16(key string) (int16, error) {
	if v, err := s.GetUInt16(key); err != nil {
		return 0, err
	} else {
		return int16(v), nil
	}
}

func (s *CacheStore) SetInt16(key string, value int16, exp time.Duration) error {
	return s.SetUInt16(key, uint16(value), exp)
}

// int32
func (s *CacheStore) GetInt32(key string) (int32, error) {
	if v, err := s.GetUInt32(key); err != nil {
		return 0, err
	} else {
		return int32(v), nil
	}
}

func (s *CacheStore) SetInt32(key string, value int32, exp time.Duration) error {
	return s.SetUInt32(key, uint32(value), exp)
}

// int64
func (s *CacheStore) GetInt64(key string) (int64, error) {
	if v, err := s.GetUInt64(key); err != nil {
		return 0, err
	} else {
		return int64(v), nil
	}
}

func (s *CacheStore) SetInt64(key string, value int64, exp time.Duration) error {
	return s.SetUInt64(key, uint64(value), exp)
}

// float32
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

// time
func (s *CacheStore) GetTime(key string) (time.Time, error) {
	var t time.Time
	data, err := s.Get(key)
	if err != nil {
		return t, err
	}
	if len(data) == 0 {
		return t, ErrNoDataForKey(key)
	}
	err = t.UnmarshalBinary(data)
	return t, err
}

func (s *CacheStore) SetTime(key string, value time.Time, exp time.Duration) error {
	if b, err := value.MarshalBinary(); err != nil {
		return err
	} else {
		return s.Set(key, b, exp)
	}
}

// JSON
func (s *CacheStore) GetJSON(key string, target interface{}) error {
	data, err := s.Get(key)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return ErrNoDataForKey(key)
	}
	return json.Unmarshal(data, target)
}

func (s *CacheStore) SetJSON(key string, value interface{}, exp time.Duration) error {
	if data, err := json.Marshal(value); err != nil {
		return err
	} else {
		return s.Set(key, data, exp)
	}
}
