package store

import (
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
)

type RWTransactionFunc func(tx RWTransaction) error

type RWTransaction interface {
	ReadTransaction
	Set(key string, entry entry.Entry) error
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
