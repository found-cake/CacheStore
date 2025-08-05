package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsRaw(key string) ([]byte, error) {
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	result := make([]byte, len(e.Data))
	copy(result, e.Data)

	return result, nil
}

func (e *Entry) AsRawNoCopy(key string) ([]byte, error) {
	if e.Type != types.RAW {
		return nil, errors.ErrTypeMismatch(types.RAW, e.Type)
	}

	return e.Data, nil
}

func FromRaw(value []byte, exp time.Duration) Entry {
	return NewEntry(types.RAW, value, exp)
}
