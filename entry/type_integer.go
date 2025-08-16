package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsInt16() (int16, error) {
	if e.Type != types.INT16 {
		return 0, errors.ErrTypeMismatch(types.INT16, e.Type)
	}
	return utils.Binary2Int16(e.Data)
}

func FromInt16(value int16, exp time.Duration) Entry {
	return NewEntry(types.INT16, utils.Int16toBinary(value), exp)
}

func (e *Entry) AsInt32() (int32, error) {
	if e.Type != types.INT32 {
		return 0, errors.ErrTypeMismatch(types.INT32, e.Type)
	}
	return utils.Binary2Int32(e.Data)
}

func FromInt32(key string, value int32, exp time.Duration) Entry {
	return NewEntry(types.INT32, utils.Int32toBinary(value), exp)
}

func (e *Entry) AsInt64() (int64, error) {
	if e.Type != types.INT64 {
		return 0, errors.ErrTypeMismatch(types.INT64, e.Type)
	}
	return utils.Binary2Int64(e.Data)
}

func FromInt64(key string, value int64, exp time.Duration) Entry {
	return NewEntry(types.INT64, utils.Int64toBinary(value), exp)
}
