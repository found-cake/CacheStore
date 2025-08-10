package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsFloat32() (float32, error) {
	if e.Type != types.FLOAT32 {
		return 0, errors.ErrTypeMismatch(types.FLOAT32, e.Type)
	}
	return utils.Binary2Float32(e.Data)
}

func FromFloat32(value float32, exp time.Duration) Entry {
	return NewEntry(types.FLOAT32, utils.Float32toBinary(value), exp)
}

func (e *Entry) AsFloat64() (float64, error) {
	if e.Type != types.FLOAT64 {
		return 0, errors.ErrTypeMismatch(types.FLOAT64, e.Type)
	}
	return utils.Binary2Float64(e.Data)
}

func FromFloat64(key string, value float64, exp time.Duration) Entry {
	return NewEntry(types.FLOAT64, utils.Float64toBinary(value), exp)
}
