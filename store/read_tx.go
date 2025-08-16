package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
)

type ReadTransactionFunc func(tx ReadTransaction) error

type ReadTransaction interface {
	Get(key string) (*entry.Entry, error)
	Exists(keys ...string) int
	TTL(key string) time.Duration
}

func (s *CacheStore) ReadTransaction(useSnapshot bool, fx ReadTransactionFunc) error {
	if s.IsClosed() {
		return errors.ErrIsClosed
	}

	if useSnapshot {
		return s.snapshotReadTx(fx)
	} else {
		return s.lockReadTx(fx)
	}
}
