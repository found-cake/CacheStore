package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type LockRWTransaction struct {
	*LockReadTransaction
	parent *CacheStore
}

func (s *CacheStore) lockRWTx(fn RWTransactionFunc) error {
	tx := &LockRWTransaction{
		parent:              s,
		LockReadTransaction: &LockReadTransaction{parent: s},
	}

	s.persistentMux.Lock()
	s.temporaryMux.Lock()
	if s.dirty != nil {
		s.dirty.mux.Lock()
	}
	defer tx.commit()

	return fn(tx)
}

func (tx *LockRWTransaction) commit() error {
	tx.parent.persistentMux.Unlock()
	tx.parent.temporaryMux.Unlock()
	if tx.parent.dirty != nil {
		tx.parent.dirty.mux.Unlock()
	}
	return nil
}

func (tx *LockRWTransaction) Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if value == nil {
		return errors.ErrValueNil
	}

	if expiry <= 0 {
		tx.parent.memorydbPersistent[key] = entry.NewEntry(dataType, value, 0)
		delete(tx.parent.memorydbTemporary, key)
	} else {
		tx.parent.memorydbTemporary[key] = entry.NewEntry(dataType, value, expiry)
		delete(tx.parent.memorydbPersistent, key)
	}

	if tx.parent.dirty != nil {
		tx.parent.dirty.set(key)
	}

	return nil
}

func (tx *LockRWTransaction) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	delete(tx.parent.memorydbTemporary, key)
	delete(tx.parent.memorydbPersistent, key)

	if tx.parent.dirty != nil {
		tx.parent.dirty.delete(key)
	}

	return nil
}
