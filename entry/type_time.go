package entry

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsTime() (time.Time, error) {
	var t time.Time
	if e.Type != types.TIME {
		return t, errors.ErrTypeMismatch(types.TIME, e.Type)
	}
	if len(e.Data) == 0 {
		return t, errors.ErrDataEmpty
	}
	err := t.UnmarshalBinary(e.Data)
	return t, err
}

func FromTime(value time.Time, exp time.Duration) (Entry, error) {
	if b, err := value.MarshalBinary(); err != nil {
		return Entry{}, err
	} else {
		return NewEntry(types.TIME, b, exp), nil
	}
}
