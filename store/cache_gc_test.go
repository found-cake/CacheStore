package store

import (
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func TestCleanExpired(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	key := "foo"
	err = store.Set(key, types.RAW, []byte("bar"), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	count := store.Exists(key)
	if count != 1 {
		t.Errorf("Exists() = %v, want 1 for expired key", count)
	}

	time.Sleep(200 * time.Millisecond)

	count = store.Exists(key)
	if count != 0 {
		t.Errorf("Exists() = %v, want 0 for expired key", count)
	}

	_, ok := store.memorydb[key]
	if !ok {
		t.Error("Want it to exist because haven't called cleanExpired yet.")
	}

	store.cleanExpired()
	_, ok = store.memorydb[key]
	if ok {
		t.Error("should not exist because cleanExpired was called.")
	}
}

func TestGarbageCollector(t *testing.T) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 500 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	key := "foo"
	err = store.Set(key, types.RAW, []byte("bar"), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	count := store.Exists(key)
	if count != 1 {
		t.Errorf("Exists() = %v, want 1 for expired key", count)
	}

	timeout := time.After(2 * time.Second)
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-timeout:
			t.Fatal("timeout: key still exists after expected GC interval")
		case <-tick:
			if _, ok := store.memorydb[key]; !ok {
				return
			}
		}
	}
}
