package errors

import (
	"errors"
	"fmt"

	"github.com/found-cake/CacheStore/utils/types"
)

type unsigned interface {
	~uint16 | ~uint32 | ~uint64
}

var (
	ErrKeyEmpty            = errors.New("key cannot be empty")
	ErrValueNil            = errors.New("value cannot be null")
	ErrDBNotInit           = errors.New("database not initialized")
	ErrFileNameEmpty       = errors.New("filename cannot be empty")
	ErrAlreadySave         = errors.New("save operation already in progress")
	ErrDirtyThresholdCount = errors.New("DirtyThresholdCount is greater than '0'")
	ErrDirtyThresholdRatio = errors.New("DirtyThresholdRatio is '0 ~ 1'")
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

func ErrUnsignedUnderflow[T unsigned](key string, current T, delta T) error {
	return fmt.Errorf("unsigned integer underflow for key '%s': current value %v is less than delta %v", key, current, delta)
}
