package store

import (
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type SnapshotReadTransaction struct {
	memorydbPersistent map[string]entry.Entry
	memorydbTemporary  map[string]entry.Entry
}

func (s *CacheStore) snapshotReadTx(fn ReadTransactionFunc) error {
	s.persistentMux.RLock()
	s.temporaryMux.RLock()

	tx := &SnapshotReadTransaction{
		memorydbPersistent: make(map[string]entry.Entry, len(s.memorydbPersistent)),
		memorydbTemporary:  make(map[string]entry.Entry, len(s.memorydbTemporary)),
	}

	for key, e := range s.memorydbPersistent {
		dataCopy := make([]byte, len(e.Data))
		copy(dataCopy, e.Data)

		tx.memorydbPersistent[key] = entry.Entry{
			Type:   e.Type,
			Data:   dataCopy,
			Expiry: e.Expiry,
		}
	}

	for key, e := range s.memorydbTemporary {
		dataCopy := make([]byte, len(e.Data))
		copy(dataCopy, e.Data)

		tx.memorydbTemporary[key] = entry.Entry{
			Type:   e.Type,
			Data:   dataCopy,
			Expiry: e.Expiry,
		}
	}

	s.persistentMux.RUnlock()
	s.temporaryMux.RUnlock()

	return fn(tx)
}

func (tx *SnapshotReadTransaction) Get(key string) (types.DataType, []byte, error) {
	t, data, err := tx.GetNoCopy(key)
	if err != nil {
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

	entry, ok := tx.memorydbPersistent[key]
	if ok {
		return entry.Type, entry.Data, nil
	}

	entry, ok = tx.memorydbTemporary[key]
	if ok {
		if entry.IsExpired() {
			return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
		} else {
			return entry.Type, entry.Data, nil
		}
	}

	return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
}

func (tx *SnapshotReadTransaction) Exists(keys ...string) int {
	if len(keys) == 0 {
		return 0
	}

	count := 0
	now := time.Now().UnixMilli()

	for _, key := range keys {
		if _, exists := tx.memorydbPersistent[key]; exists {
			count++
		} else if entry, exists := tx.memorydbTemporary[key]; exists {
			if !entry.IsExpiredWithUnixMilli(now) {
				count++
			}
		}
	}

	return count
}

func (tx *SnapshotReadTransaction) TTL(key string) time.Duration {
	_, ok := tx.memorydbPersistent[key]
	if ok {
		return TTLNoExpiry
	}

	entry, ok := tx.memorydbTemporary[key]
	if !ok {
		return TTLExpired
	}
	now := time.Now().UnixMilli()
	if now >= entry.Expiry {
		return TTLExpired
	}

	remaining := time.Duration(entry.Expiry-now) * time.Millisecond
	return remaining
}
