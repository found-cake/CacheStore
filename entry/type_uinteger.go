package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsUInt16() (uint16, error) {
	if e.Type != types.UINT16 {
		return 0, errors.ErrTypeMismatch(types.UINT16, e.Type)
	}
	return utils.Binary2UInt16(e.Data)
}

func FromUInt16(value uint16, exp time.Duration) Entry {
	return NewEntry(types.UINT16, utils.UInt16toBinary(value), exp)
}

func (e *Entry) AsUInt32() (uint32, error) {
	if e.Type != types.UINT32 {
		return 0, errors.ErrTypeMismatch(types.UINT32, e.Type)
	}
	return utils.Binary2UInt32(e.Data)
}

func FromUInt32(value uint32, exp time.Duration) Entry {
	return NewEntry(types.UINT32, utils.UInt32toBinary(value), exp)
}

func (e *Entry) AsUInt64() (uint64, error) {
	if e.Type != types.UINT64 {
		return 0, errors.ErrTypeMismatch(types.UINT64, e.Type)
	}
	return utils.Binary2UInt64(e.Data)
}

func FromUInt64(value uint64, exp time.Duration) Entry {
	return NewEntry(types.UINT64, utils.UInt64toBinary(value), exp)
}
