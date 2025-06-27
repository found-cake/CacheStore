package errors

import (
	"errors"
	"fmt"

	"github.com/found-cake/CacheStore/store/types"
)

var (
	ErrKeyEmpty      = errors.New("key cannot be empty")
	ErrValueNil      = errors.New("value cannot be null")
	ErrDBNotInit     = errors.New("database not initialized")
	ErrFileNameEmpty = errors.New("filename cannot be empty")
)

func ErrInvalidDataLength(expected, actual int) error {
	return fmt.Errorf("invalid data length: expected %d bytes, got %d bytes", expected, actual)
}

func ErrNoDataForKey(key string) error {
	return fmt.Errorf("no data found for key: %s", key)
}

func ErrTypeMismatch(key string, expected, actual types.DataType) error {
	return fmt.Errorf("type mismatch for key '%s': expected %s, got %s",
		key, expected.String(), actual.String())
}
