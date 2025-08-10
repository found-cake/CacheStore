package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils"
	"github.com/found-cake/CacheStore/utils/generic"
	"github.com/found-cake/CacheStore/utils/types"
)

func (s *CacheStore) GetUInt16(key string) (uint16, error) {
	return s.getNum16(key, types.UINT16)
}

func (s *CacheStore) SetUInt16(key string, value uint16, exp time.Duration) error {
	return s.Set(key, types.UINT16, utils.UInt16toBinary(value), exp)
}

func (s *CacheStore) IncrUInt16(key string, delta uint16, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.UINT16, exp,
		utils.Binary2UInt16,
		utils.UInt16toBinary,
		utils.UInt16CheckOverFlow,
		nil,
	)
}

func (s *CacheStore) DecrUInt16(key string, delta uint16, exp time.Duration) error {
	return decrUnsigned(
		s, key, delta, types.UINT16, exp,
		utils.Binary2UInt16,
		utils.UInt16toBinary,
		utils.UintCheckUnderFlow,
	)
}

func (s *CacheStore) GetUInt32(key string) (uint32, error) {
	return s.getNum32(key, types.UINT32)
}

func (s *CacheStore) SetUInt32(key string, value uint32, exp time.Duration) error {
	return s.Set(key, types.UINT32, utils.UInt32toBinary(value), exp)
}

func (s *CacheStore) IncrUInt32(key string, delta uint32, exp time.Duration) error {
	return incrNumber(
		s, key, delta, types.UINT32, exp,
		utils.Binary2UInt32,
		utils.UInt32toBinary,
		utils.UInt32CheckOverFlow,
		nil,
	)
}

func (s *CacheStore) DecrUInt32(key string, delta uint32, exp time.Duration) error {
	return decrUnsigned(
		s, key, delta, types.UINT32, exp,
		utils.Binary2UInt32,
		utils.UInt32toBinary,
		utils.UintCheckUnderFlow,
	)
}

func (s *CacheStore) GetUInt64(key string) (uint64, error) {
	return s.getNum64(key, types.UINT64)
}

func (s *CacheStore) SetUInt64(key string, value uint64, exp time.Duration) error {
	return s.Set(key, types.UINT64, utils.UInt64toBinary(value), exp)
}

func (s *CacheStore) IncrUInt64(key string, delta uint64, exp time.Duration) error {
	return incrNumber(s, key, delta, types.UINT64, exp, utils.Binary2UInt64, utils.UInt64toBinary, utils.UInt64CheckOverFlow, nil)
}

func (s *CacheStore) DecrUInt64(key string, delta uint64, exp time.Duration) error {
	return decrUnsigned(
		s, key, delta, types.UINT64, exp,
		utils.Binary2UInt64,
		utils.UInt64toBinary,
		utils.UintCheckUnderFlow,
	)
}

func decrUnsigned[T generic.Unsigned](
	s *CacheStore,
	key string,
	delta T,
	data_type types.DataType,
	exp time.Duration,
	fromBinary func([]byte) (T, error),
	toBinary func(T) []byte,
	checkUnderflow func(T, T) bool,
) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	return s.RWTransaction(false, func(tx RWTransaction) error {
		e, err := tx.Get(key)
		if err != nil {
			return errors.ErrNoDataForKey(key)
		}
		if e.Type != data_type {
			return errors.ErrTypeMismatch(data_type, e.Type)
		}
		value, err := fromBinary(e.Data)
		if err != nil {
			return err
		}
		if checkUnderflow(value, delta) {
			return errors.ErrUnsignedUnderflow(key, value, delta)
		}
		value -= delta
		data := toBinary(value)
		tx.Set(key, entry.NewEntry(data_type, data, exp))
		return nil
	})
}
