package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type RWTransactionFunc func(tx RWTransaction) error

type RWTransaction interface {
	ReadTransaction
	Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error
	Delete(key string) error
	commit() error
}

func (s *CacheStore) RWTransaction(useSnapshot bool, fx RWTransactionFunc) error {
	if s.IsClosed() {
		return errors.ErrIsClosed
	}

	if useSnapshot {
		return s.snapshotRwTx(fx)
	} else {
		return s.lockRWTx(fx)
	}
}
