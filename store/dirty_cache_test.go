package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func TestNewDirtyManager(t *testing.T) {
	count := 100
	ratio := 0.5

	dm := newDirtyManager(count, ratio)

	if dm.ThresholdCount != count {
		t.Errorf("ThresholdCount = %d, want %d", dm.ThresholdCount, count)
	}
	if dm.ThresholdRatio != ratio {
		t.Errorf("ThresholdRatio = %f, want %f", dm.ThresholdRatio, ratio)
	}
	if dm.needFullSync != false {
		t.Errorf("needFullSync = %t, want false", dm.needFullSync)
	}
	if len(dm.dirtyData) != 0 {
		t.Errorf("dirtyData length = %d, want 0", len(dm.dirtyData))
	}
}

func TestDirtyManager_Set(t *testing.T) {
	dm := newDirtyManager(10, 0.5)

	dm.set("key1")
	if dm.size() != 1 {
		t.Errorf("size after set = %d, want 1", dm.size())
	}

	dm.unsafeSet("key2")
	if dm.size() != 2 {
		t.Errorf("size after unsafeSet = %d, want 2", dm.size())
	}

	dm.set("key1")
	if dm.size() != 2 {
		t.Errorf("size after duplicate set = %d, want 2", dm.size())
	}

	setKeys, deleteKeys := dm.keys()
	if len(setKeys) != 2 {
		t.Errorf("setKeys length = %d, want 2", len(setKeys))
	}
	if len(deleteKeys) != 0 {
		t.Errorf("deleteKeys length = %d, want 0", len(deleteKeys))
	}

	expectedKeys := map[string]bool{"key1": true, "key2": true}
	for _, key := range setKeys {
		if !expectedKeys[key] {
			t.Errorf("unexpected set key: %s", key)
		}
	}
}

func TestDirtyManager_Delete(t *testing.T) {
	dm := newDirtyManager(10, 0.5)

	dm.delete("key1")
	if dm.size() != 1 {
		t.Errorf("size after delete = %d, want 1", dm.size())
	}

	dm.unsafeDelete("key2")
	if dm.size() != 2 {
		t.Errorf("size after unsafeDelete = %d, want 2", dm.size())
	}

	dm.delete("key1")
	if dm.size() != 2 {
		t.Errorf("size after duplicate delete = %d, want 2", dm.size())
	}

	setKeys, deleteKeys := dm.keys()
	if len(setKeys) != 0 {
		t.Errorf("setKeys length = %d, want 0", len(setKeys))
	}
	if len(deleteKeys) != 2 {
		t.Errorf("deleteKeys length = %d, want 2", len(deleteKeys))
	}

	expectedKeys := map[string]bool{"key1": true, "key2": true}
	for _, key := range deleteKeys {
		if !expectedKeys[key] {
			t.Errorf("unexpected delete key: %s", key)
		}
	}
}

func TestDirtyManager_MixedOperations(t *testing.T) {
	dm := newDirtyManager(10, 0.5)

	dm.set("key1")
	dm.set("key2")
	dm.set("key3")

	dm.delete("key4")
	dm.delete("key5")

	dm.delete("key1")

	dm.set("key4")

	if dm.size() != 5 {
		t.Errorf("size = %d, want 5", dm.size())
	}

	setKeys, deleteKeys := dm.keys()
	if len(setKeys) != 3 {
		t.Errorf("setKeys length = %d, want 3", len(setKeys))
	}
	if len(deleteKeys) != 2 {
		t.Errorf("deleteKeys length = %d, want 2", len(deleteKeys))
	}

	expectedSetKeys := map[string]bool{"key2": true, "key3": true, "key4": true}
	for _, key := range setKeys {
		if !expectedSetKeys[key] {
			t.Errorf("unexpected set key: %s", key)
		}
	}

	expectedDeleteKeys := map[string]bool{"key1": true, "key5": true}
	for _, key := range deleteKeys {
		if !expectedDeleteKeys[key] {
			t.Errorf("unexpected delete key: %s", key)
		}
	}
}

func TestDirtyManager_Clear(t *testing.T) {
	dm := newDirtyManager(10, 0.5)

	dm.set("key1")
	dm.delete("key2")
	dm.set("key3")

	if dm.size() != 3 {
		t.Errorf("size before clear = %d, want 3", dm.size())
	}

	dm.clear()
	if dm.size() != 0 {
		t.Errorf("size after clear = %d, want 0", dm.size())
	}

	dm.set("key4")
	dm.delete("key5")

	dm.unsafeClear()
	if dm.size() != 0 {
		t.Errorf("size after unsafeClear = %d, want 0", dm.size())
	}

	setKeys, deleteKeys := dm.keys()
	if len(setKeys) != 0 {
		t.Errorf("setKeys length after clear = %d, want 0", len(setKeys))
	}
	if len(deleteKeys) != 0 {
		t.Errorf("deleteKeys length after clear = %d, want 0", len(deleteKeys))
	}
}

func TestDirtyManager_WantFullSync(t *testing.T) {
	dm := newDirtyManager(10, 0.5)

	dm.set("key1")
	dm.delete("key2")
	dm.set("key3")

	if dm.size() != 3 {
		t.Errorf("size before NeedFullSync = %d, want 3", dm.size())
	}
	if dm.needFullSync != false {
		t.Errorf("needFullSync before call = %t, want false", dm.needFullSync)
	}

	dm.wantFullSync()

	if dm.needFullSync != true {
		t.Errorf("needFullSync after call = %t, want true", dm.needFullSync)
	}
	if dm.size() != 0 {
		t.Errorf("size after NeedFullSync = %d, want 0", dm.size())
	}

	dm.set("key4")
	dm.delete("key5")

	if dm.size() != 0 {
		t.Errorf("size after operations with needFullSync = %d, want 0", dm.size())
	}

	setKeys, deleteKeys := dm.keys()
	if len(setKeys) != 0 {
		t.Errorf("setKeys length with needFullSync = %d, want 0", len(setKeys))
	}
	if len(deleteKeys) != 0 {
		t.Errorf("deleteKeys length with needFullSync = %d, want 0", len(deleteKeys))
	}
}

func TestDirtyManager_WithCacheStore(t *testing.T) {
	path := filepath.Join(t.TempDir(), "testdb.sqlite")

	store, err := NewCacheStore(config.Config{
		DBSave:              true,
		DBFileName:          path,
		DBSaveInterval:      time.Hour,
		SaveDirtyData:       true,
		DirtyThresholdCount: 5,
		DirtyThresholdRatio: 0.5,
	})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if store.dirty == nil {
		t.Fatal("dirty manager should be initialized")
	}

	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err := store.Set(key, types.STRING, []byte("value"), time.Hour)
		if err != nil {
			t.Errorf("Set() error = %v", err)
		}
	}

	if store.dirty.size() != len(keys) {
		t.Errorf("dirty manager size = %d, want %d", store.dirty.size(), len(keys))
	}

	err = store.Delete("key1")
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	if store.dirty.size() != 3 {
		t.Errorf("dirty manager size after delete = %d, want 3", store.dirty.size())
	}

	setKeys, deleteKeys := store.dirty.keys()
	if len(setKeys) != 2 {
		t.Errorf("setKeys length = %d, want 2", len(setKeys))
	}
	if len(deleteKeys) != 1 {
		t.Errorf("deleteKeys length = %d, want 1", len(deleteKeys))
	}

	store.Flush()

	if !store.dirty.needFullSync {
		t.Error("needFullSync should be true after Flush")
	}

	if store.dirty.size() != 0 {
		t.Errorf("dirty manager size after Flush = %d, want 0", store.dirty.size())
	}
}
