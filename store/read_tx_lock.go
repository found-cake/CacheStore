package store

import (
	"time"

	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

type LockReadTransaction struct {
	parent *CacheStore
}

func (s *CacheStore) lockReadTx(fn ReadTransactionFunc) error {
	tx := &LockReadTransaction{
		parent: s,
	}

	s.persistentMux.RLock()
	s.temporaryMux.RLock()
	defer func() {
		s.persistentMux.RUnlock()
		s.temporaryMux.RUnlock()
	}()

	return fn(tx)
}

func (tx *LockReadTransaction) Get(key string) (types.DataType, []byte, error) {
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
func (tx *LockReadTransaction) GetNoCopy(key string) (types.DataType, []byte, error) {
	if key == "" {
		return types.UNKNOWN, nil, errors.ErrKeyEmpty
	}

	entry, ok := tx.parent.memorydbPersistent[key]
	if ok {
		return entry.Type, entry.Data, nil
	}

	entry, ok = tx.parent.memorydbTemporary[key]
	if ok {
		if entry.IsExpired() {
			return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
		} else {
			return entry.Type, entry.Data, nil
		}
	}

	return types.UNKNOWN, nil, errors.ErrNoDataForKey(key)
}

func (tx *LockReadTransaction) Exists(keys ...string) int {
	if len(keys) == 0 {
		return 0
	}

	count := 0
	now := time.Now().UnixMilli()

	for _, key := range keys {
		if _, exists := tx.parent.memorydbPersistent[key]; exists {
			count++
		} else if entry, exists := tx.parent.memorydbTemporary[key]; exists {
			if !entry.IsExpiredWithUnixMilli(now) {
				count++
			}
		}
	}

	return count
}

func (tx *LockReadTransaction) TTL(key string) time.Duration {
	_, ok := tx.parent.memorydbPersistent[key]
	if ok {
		return TTLNoExpiry
	}

	entry, ok := tx.parent.memorydbTemporary[key]
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
