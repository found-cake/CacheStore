package entry

import (
	"time"

	"github.com/found-cake/CacheStore/utils/types"
)

type Entry struct {
	Type   types.DataType
	Data   []byte
	Expiry int64
}

func (e Entry) IsExpired() bool {
	return e.Expiry > 0 && e.Expiry <= time.Now().Unix()
}

func (e Entry) IsExpiredWithTime(now int64) bool {
	return e.Expiry > 0 && e.Expiry <= now
}

func NewEntry(dataType types.DataType, data []byte, exp time.Duration) Entry {
	var expiry int64
	if exp > 0 {
		expiry = time.Now().Add(exp).Unix()
	}
	return Entry{
		Type:   dataType,
		Data:   data,
		Expiry: expiry,
	}
}
