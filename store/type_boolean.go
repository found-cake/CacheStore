package store

import "time"

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
