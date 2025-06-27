package errors

import (
	"errors"
	"fmt"
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
