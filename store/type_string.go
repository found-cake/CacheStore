package store

import "time"

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
