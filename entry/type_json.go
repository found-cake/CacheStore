package entry

import (
	"encoding/json"
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func (e *Entry) AsJSON(target interface{}) error {
	if e.Type != types.JSON {
		return errors.ErrTypeMismatch(types.JSON, e.Type)
	}
	if len(e.Data) == 0 {
		return errors.ErrDataEmpty
	}
	return json.Unmarshal(e.Data, target)
}

func FromJSON(value interface{}, exp time.Duration) (Entry, error) {
	if data, err := json.Marshal(value); err != nil {
		return Entry{}, err
	} else {
		return NewEntry(types.JSON, data, exp), nil
	}
}
