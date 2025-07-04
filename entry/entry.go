package entry

import (
	"time"

	"github.com/found-cake/CacheStore/utils/types"
)

type Entry struct {
	Type   types.DataType
	Data   []byte
	Expiry uint32
}

func (e Entry) IsExpired() bool {
	return e.Expiry > 0 && e.Expiry <= uint32(time.Now().Unix())
}

func (e Entry) IsExpiredWithTime(now uint32) bool {
	return e.Expiry > 0 && e.Expiry <= now
}

func NewEntry(dataType types.DataType, data []byte, exp time.Duration) Entry {
	var expiry uint32
	if exp > 0 {
		expiry = uint32(time.Now().Add(exp).Unix())
	}
	return Entry{
		Type:   dataType,
		Data:   data,
		Expiry: expiry,
	}
}
