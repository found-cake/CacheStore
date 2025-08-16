package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsBool() (bool, error) {
	if e.Type != types.BOOLEAN {
		return false, errors.ErrTypeMismatch(types.BOOLEAN, e.Type)
	}
	return len(e.Data) > 0 && e.Data[0] == 1, nil
}

func FromBool(value bool, exp time.Duration) Entry {
	v := byte(0)
	if value {
		v = 1
	}
	return NewEntry(types.BOOLEAN, []byte{v}, exp)
}
