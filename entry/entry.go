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

func (e *Entry) CopyData() []byte {
	result := make([]byte, len(e.Data))
	copy(result, e.Data)
	return result
}

func (e *Entry) IsExpired() bool {
	return e.IsExpiredWithUnixMilli(time.Now().UnixMilli())
}

func (e *Entry) IsExpiredWithUnixMilli(now int64) bool {
	return e.Expiry > 0 && e.Expiry <= now
}

func NewEntry(dataType types.DataType, data []byte, exp time.Duration) Entry {
	var expiry int64
	if exp > 0 {
		expiry = time.Now().Add(exp).UnixMilli()
	}
	return Entry{
		Type:   dataType,
		Data:   data,
		Expiry: expiry,
	}
}
