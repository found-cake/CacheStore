package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsRaw() ([]byte, error) {
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	return e.CopyData(), nil
}

func (e *Entry) AsRawNoCopy() ([]byte, error) {
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	return e.Data, nil
}

func FromRaw(value []byte, exp time.Duration) Entry {
	return NewEntry(types.RAW, value, exp)
}
