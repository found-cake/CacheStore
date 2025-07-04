package errors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/found-cake/CacheStore/utils/generic"
	"github.com/found-cake/CacheStore/utils/types"
)

var (
	ErrKeyEmpty            = errors.New("key cannot be empty")
	ErrValueNil            = errors.New("value cannot be null")
	ErrDBNotInit           = errors.New("database not initialized")
	ErrFileNameEmpty       = errors.New("filename cannot be empty")
	ErrAlreadySave         = errors.New("save operation already in progress")
	ErrDirtyThresholdCount = errors.New("DirtyThresholdCount is greater than '0'")
	ErrDirtyThresholdRatio = errors.New("DirtyThresholdRatio is '0 ~ 1'")
	ErrFloatSpecial        = errors.New("Invalid Error: result is Nan(Not a Number) or Infinity")
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

func ErrUnsignedUnderflow[T generic.Unsigned](key string, current, delta T) error {
	return fmt.Errorf("unsigned integer underflow for key '%s': current value %v is less than delta %v", key, current, delta)
}

func ErrValueOverflow[T generic.Numberic](key string, data_type types.DataType, current, delta T) error {
	return fmt.Errorf("%s overflow for key '%s': %v + %v exceeds representable range", strings.ToLower(data_type.String()), key, current, delta)
}
