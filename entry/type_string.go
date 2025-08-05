package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsString() (string, error) {
	if e.Type != types.STRING {
		return "", errors.ErrTypeMismatch(types.STRING, e.Type)
	}
	return string(e.Data), nil
}

func FromString(value string, exp time.Duration) Entry {
	return NewEntry(types.STRING, []byte(value), exp)
}
