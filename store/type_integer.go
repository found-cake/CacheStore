package store

import "time"

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
