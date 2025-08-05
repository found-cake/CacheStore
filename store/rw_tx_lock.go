package store

import (
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
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

func (tx *LockRWTransaction) Set(key string, e entry.Entry) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if e.Data == nil {
		return errors.ErrValueNil
	}

	if e.Expiry <= 0 {
		tx.parent.memorydbPersistent[key] = e
		delete(tx.parent.memorydbTemporary, key)
	} else {
		tx.parent.memorydbTemporary[key] = e
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
