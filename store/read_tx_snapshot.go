package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type SnapshotReadTransaction struct {
	memorydb map[string]entry.Entry
}

func newSnapshotReadTX(s *CacheStore) *SnapshotReadTransaction {
	s.persistentMux.RLock()
	s.temporaryMux.RLock()

	tx := &SnapshotReadTransaction{
		memorydb: make(map[string]entry.Entry, len(s.memorydbPersistent)+len(s.memorydbTemporary)),
	}

	for key, e := range s.memorydbPersistent {
		dataCopy := make([]byte, len(e.Data))
		copy(dataCopy, e.Data)

		tx.memorydb[key] = entry.Entry{
			Type:   e.Type,
			Data:   dataCopy,
			Expiry: e.Expiry,
		}
	}

	for key, e := range s.memorydbTemporary {
		if e.IsExpired() {
			continue
		}
		dataCopy := make([]byte, len(e.Data))
		copy(dataCopy, e.Data)

		tx.memorydb[key] = entry.Entry{
			Type:   e.Type,
			Data:   dataCopy,
			Expiry: e.Expiry,
		}
	}

	s.persistentMux.RUnlock()
	s.temporaryMux.RUnlock()

	return tx
}

func (s *CacheStore) snapshotReadTx(fn ReadTransactionFunc) error {
	tx := newSnapshotReadTX(s)
	return fn(tx)
}

func (tx *SnapshotReadTransaction) Get(key string) (types.DataType, []byte, error) {
	t, data, err := tx.GetNoCopy(key)
	if err == nil {
		result := make([]byte, len(data))
		copy(result, data)
		return t, result, err
	}
	return t, data, err
}

// GetNoCopy retrieves a value without copying data (zero-copy read)
// ⚠️ WARNING: Don't modify the returned value!
func (tx *SnapshotReadTransaction) GetNoCopy(key string) (types.DataType, []byte, error) {
	if key == "" {
		return types.UNKNOWN, nil, errors.ErrKeyEmpty
	}

	entry, ok := tx.memorydb[key]
	if !ok {
		return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
	}
	if entry.IsExpired() {
		return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
	}

	return entry.Type, entry.Data, nil
}

func (tx *SnapshotReadTransaction) Exists(keys ...string) int {
	if len(keys) == 0 {
		return 0
	}

	count := 0
	now := time.Now().UnixMilli()

	for _, key := range keys {
		if entry, exists := tx.memorydb[key]; exists {
			if !entry.IsExpiredWithUnixMilli(now) {
				count++
			}
		}
	}

	return count
}

func (tx *SnapshotReadTransaction) TTL(key string) time.Duration {
	e, ok := tx.memorydb[key]
	if !ok {
		return TTLExpired
	}

	if e.Expiry == 0 {
		return TTLNoExpiry
	}

	now := time.Now().UnixMilli()
	if now >= e.Expiry {
		return TTLExpired
	}

	remaining := time.Duration(e.Expiry-now) * time.Millisecond
	return remaining
}
