package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/generic"
	"github.com/found-cake/CacheStore/utils/types"
)

func incrNumber[T generic.Numberic](
	s *CacheStore,
	key string,
	delta T,
	data_type types.DataType,
	exp time.Duration,
	fromBinary func([]byte) (T, error),
	toBinary func(T) []byte,
	checkOverFlow func(T, T) bool,
	checkFloatSpesial func(T) bool,
) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	return s.RWTransaction(false, func(tx RWTransaction) error {
		e, err := tx.Get(key)
		if err != nil {
			data := toBinary(delta)
			tx.Set(key, entry.NewEntry(data_type, data, exp))
			return nil
		}
		if e.Type != data_type {
			return errors.ErrTypeMismatch(data_type, e.Type)
		}
		value, err := fromBinary(e.Data)
		if err != nil {
			return err
		}
		if checkOverFlow(value, delta) {
			return errors.ErrValueOverflow(key, data_type, value, delta)
		}
		value += delta
		data := toBinary(value)
		if checkFloatSpesial != nil && checkFloatSpesial(value) {
			return errors.ErrFloatSpecial
		}
		tx.Set(key, entry.NewEntry(data_type, data, exp))
		return nil
	})
}
