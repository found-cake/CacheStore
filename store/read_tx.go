package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type ReadTransactionFunc func(tx ReadTransaction) error

type ReadTransaction interface {
	Get(key string) (types.DataType, []byte, error)
	GetNoCopy(key string) (types.DataType, []byte, error)
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
