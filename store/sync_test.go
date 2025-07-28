package store

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func TestSync_Basic(t *testing.T) {
	dbFile := tempDBFile(t)
	defer os.Remove(dbFile)

	store, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dbFile,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	testData := map[string]string{
		"user:1001": "Alice",
		"user:1002": "Bob",
		"cache:key": "some_value",
	}
	for key, value := range testData {
		err = store.Set(key, types.STRING, []byte(value), time.Hour)
		if err != nil {
			t.Fatalf("Set() failed for [%s]: %v", key, err)
		}
	}

	store.Sync()
	store.Close()

	store2, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dbFile,
	})
	if err != nil {
		t.Fatalf("Failed to re-create store: %v", err)
	}
	defer store2.Close()

	for key, want := range testData {
		dtype, got, err := store2.Get(key)
		if err != nil {
			t.Errorf("Get() failed for key [%s]: %v", key, err)
			continue
		}
		if dtype != types.STRING {
			t.Errorf("Type mismatch for [%s]: got %v, want %v", key, dtype, types.STRING)
		}
		if string(got) != want {
			t.Errorf("Value mismatch for [%s]: got %s, want %s", key, string(got), want)
		}
	}
}

func TestSync_FullSync(t *testing.T) {
	fullDB := tempDBFile(t)
	defer os.Remove(fullDB)
	store, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: fullDB,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bulk_key_%d", i)
		val := fmt.Sprintf("bulk_value_%d", i)
		if err := store.Set(key, types.STRING, []byte(val), time.Hour); err != nil {
			t.Fatalf("Set() failed: %v", err)
		}
	}
	store.FullSync()
	store.Close()

	store2, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: fullDB,
	})
	if err != nil {
		t.Fatalf("Failed to re-create store: %v", err)
	}
	defer store2.Close()
	keys := store2.Keys()
	if len(keys) != 100 {
		t.Errorf("Restored key count mismatch: got %d, want 100", len(keys))
	}
}

func TestSync_Dirty(t *testing.T) {
	dirtyDB := tempDBFile(t)
	defer os.Remove(dirtyDB)
	store, err := NewCacheStore(config.Config{
		DBSave:              true,
		DBFileName:          dirtyDB,
		SaveDirtyData:       true,
		DirtyThresholdCount: 5,
		DirtyThresholdRatio: 0.5,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("dirty_key_%d", i)
		val := fmt.Sprintf("dirty_value_%d", i)
		if err := store.Set(key, types.STRING, []byte(val), time.Hour); err != nil {
			t.Fatalf("Set() failed: %v", err)
		}
	}
	store.Sync()
	store.Close()

	store2, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dirtyDB,
	})
	if err != nil {
		t.Fatalf("Failed to re-create store: %v", err)
	}
	defer store2.Close()
	for i := 0; i < 3; i++ {
		key := fmt.Sprintf("dirty_key_%d", i)
		want := fmt.Sprintf("dirty_value_%d", i)
		dtype, got, err := store2.Get(key)
		if err != nil {
			t.Errorf("Get() failed for dirty key [%s]: %v", key, err)
			continue
		}
		if dtype != types.STRING {
			t.Errorf("Type mismatch for [%s]: got %v, want %v", key, dtype, types.STRING)
		}
		if string(got) != want {
			t.Errorf("Dirty key [%s] mismatch: got %s, want %s", key, string(got), want)
		}
	}
}

func TestSync_AutoPeriodic(t *testing.T) {
	dbFile := tempDBFile(t)
	defer os.Remove(dbFile)
	store, err := NewCacheStore(config.Config{
		DBSave:              true,
		DBFileName:          dbFile,
		SaveDirtyData:       true,
		DirtyThresholdCount: 10,
		DirtyThresholdRatio: 0.5,
		DBSaveInterval:      100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	if err := store.Set("auto_sync_key", types.STRING, []byte("auto_sync_value"), time.Hour); err != nil {
		t.Fatalf("Set() failed: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	store.Close()

	store2, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dbFile,
	})
	if err != nil {
		t.Fatalf("Failed to re-create store: %v", err)
	}
	defer store2.Close()
	_, got, err := store2.Get("auto_sync_key")
	if err != nil {
		t.Errorf("Get() failed for auto_sync_key: %v", err)
	} else if string(got) != "auto_sync_value" {
		t.Errorf("Value mismatch: got %s, want auto_sync_value", string(got))
	}
}

func TestSync_ServerRestart(t *testing.T) {
	dbFile := tempDBFile(t)
	defer os.Remove(dbFile)

	store1, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dbFile,
	})
	if err != nil {
		t.Fatalf("Failed to create first store: %v", err)
	}
	sessionData := map[string]string{
		"session:abc123": "user_id=1001",
		"session:def456": "user_id=1002",
		"config:theme":   "dark_mode",
	}
	for k, v := range sessionData {
		store1.Set(k, types.STRING, []byte(v), 24*time.Hour)
	}
	store1.Sync()
	store1.Close()

	store2, err := NewCacheStore(config.Config{
		DBSave:     true,
		DBFileName: dbFile,
	})
	if err != nil {
		t.Fatalf("Failed to create second store: %v", err)
	}
	defer store2.Close()
	for k, want := range sessionData {
		dtype, got, err := store2.Get(k)
		if err != nil {
			t.Errorf("Get() after restart failed for [%s]: %v", k, err)
			continue
		}
		if dtype != types.STRING {
			t.Errorf("Type mismatch for [%s]: got %v, want %v", k, dtype, types.STRING)
		}
		if string(got) != want {
			t.Errorf("After restart, value mismatch for [%s]: got %s, want %s", k, string(got), want)
		}
	}
	store2.Set("session:abc123", types.STRING, []byte("user_id=1001,updated=true"), time.Hour)
	store2.Sync()
}
