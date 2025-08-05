package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type SnapshotRWTransaction struct {
	*SnapshotReadTransaction
	parent            *CacheStore
	pendingPersistent map[string]*entry.Entry
	pendingTemporary  map[string]*entry.Entry
	committed         bool
}

func (s *CacheStore) snapshotRwTx(fn RWTransactionFunc) error {
	tx := &SnapshotRWTransaction{
		SnapshotReadTransaction: newSnapshotReadTX(s),
		parent:                  s,
		pendingPersistent:       make(map[string]*entry.Entry),
		pendingTemporary:        make(map[string]*entry.Entry),
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.commit()
}

func (tx *SnapshotRWTransaction) commit() error {
	if tx.committed {
		return errors.ErrAlreadyCommit
	}

	var delete_keys map[string]struct{}

	tx.parent.persistentMux.Lock()
	tx.parent.temporaryMux.Lock()
	if tx.parent.dirty != nil {
		tx.parent.dirty.mux.Lock()
		delete_keys = make(map[string]struct{}, len(tx.pendingPersistent))
		defer tx.parent.dirty.mux.Unlock()
	}
	for key, entry := range tx.pendingPersistent {
		if entry == nil {
			delete(tx.parent.memorydbPersistent, key)
			if tx.parent.dirty != nil {
				delete_keys[key] = struct{}{}
			}
		} else {
			tx.parent.memorydbPersistent[key] = *entry
			if tx.parent.dirty != nil {
				tx.parent.dirty.unsafeSet(key)
			}
		}
	}
	tx.parent.persistentMux.Unlock()

	for key, entry := range tx.pendingTemporary {
		if entry == nil {
			delete(tx.parent.memorydbTemporary, key)
		} else {
			tx.parent.memorydbTemporary[key] = *entry
			if tx.parent.dirty != nil {
				tx.parent.dirty.unsafeSet(key)
				delete(delete_keys, key)
			}
		}
	}
	tx.parent.temporaryMux.Unlock()

	if tx.parent.dirty != nil {
		for key := range delete_keys {
			tx.parent.dirty.unsafeDelete(key)
		}
	}

	tx.committed = true
	return nil
}

func (tx *SnapshotRWTransaction) Set(key string, dataType types.DataType, value []byte, expiry time.Duration) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}
	if value == nil {
		return errors.ErrValueNil
	}

	var e entry.Entry
	if expiry <= 0 {
		e = entry.NewEntry(dataType, value, 0)
		tx.pendingPersistent[key] = &e
		tx.pendingTemporary[key] = nil
	} else {
		e = entry.NewEntry(dataType, value, expiry)
		tx.pendingTemporary[key] = &e
		tx.pendingPersistent[key] = nil
	}
	tx.memorydb[key] = e

	return nil
}

func (tx *SnapshotRWTransaction) Delete(key string) error {
	if key == "" {
		return errors.ErrKeyEmpty
	}

	tx.pendingPersistent[key] = nil
	tx.pendingTemporary[key] = nil
	delete(tx.memorydb, key)

	return nil
}
