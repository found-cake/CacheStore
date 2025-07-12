package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func tempDBFile(t *testing.T) string {
	dir := t.TempDir()
	return filepath.Join(dir, "testdb.sqlite")
}

func TestCreateStore(t *testing.T) {
	tests := []struct {
		name    string
		config  config.Config
		wantErr bool
	}{
		{
			name:    "default config store",
			config:  config.DefaultConfig(),
			wantErr: false,
		},
		{
			name: "basic store without persistence",
			config: config.Config{
				DBSave: false,
			},
			wantErr: false,
		},
		{
			name: "store with empty db filename",
			config: config.Config{
				DBSave:     true,
				DBFileName: "",
			},
			wantErr: true,
		},
		{
			name: "store with invalid dirty threshold count",
			config: config.Config{
				DBSave:              true,
				SaveDirtyData:       true,
				DirtyThresholdCount: 0,
				DirtyThresholdRatio: 0.5,
				DBFileName:          tempDBFile(t),
			},
			wantErr: true,
		},
		{
			name: "store with invalid dirty threshold ratio",
			config: config.Config{
				DBSave:              true,
				SaveDirtyData:       true,
				DirtyThresholdCount: 100,
				DirtyThresholdRatio: 0,
				DBFileName:          tempDBFile(t),
			},
			wantErr: true,
		},
		{
			name: "store with valid config",
			config: config.Config{
				DBSave:              true,
				SaveDirtyData:       true,
				DirtyThresholdCount: 100,
				DirtyThresholdRatio: 0.5,
				DBFileName:          tempDBFile(t),
				DBSaveInterval:      time.Second,
				GCInterval:          time.Second,
			},
			wantErr: false,
		},
	}
	defer os.Remove(config.DefaultConfig().DBFileName)

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			store, err := NewCacheStore(tcase.config)
			if (err != nil) != tcase.wantErr {
				t.Errorf("NewCacheStore() error = %v, wantErr %v", err, tcase.wantErr)
				return
			}
			if store != nil {
				defer store.Close()
			}
		})
	}
}

func TestCacheStore_GetSet(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name     string
		key      string
		dataType types.DataType
		value    []byte
		expiry   time.Duration
		wantErr  bool
	}{
		{
			name:     "valid set and get",
			key:      "test_key",
			dataType: types.STRING,
			value:    []byte("test_value"),
			expiry:   time.Hour,
			wantErr:  false,
		},
		{
			name:     "empty key",
			key:      "",
			dataType: types.STRING,
			value:    []byte("test_value"),
			expiry:   time.Hour,
			wantErr:  true,
		},
		{
			name:     "nil value",
			key:      "test_key",
			dataType: types.STRING,
			value:    nil,
			expiry:   time.Hour,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(tt.key, tt.dataType, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gotType, gotValue, err := store.Get(tt.key)
				if err != nil {
					t.Errorf("Get() error = %v", err)
					return
				}
				if gotType != tt.dataType {
					t.Errorf("Get() gotType = %v, want %v", gotType, tt.dataType)
				}
				if string(gotValue) != string(tt.value) {
					t.Errorf("Get() gotValue = %v, want %v", string(gotValue), string(tt.value))
				}
			}
		})
	}
}

func TestCacheStore_GetNoCopy(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	key := "test_key"
	value := []byte("test_value")

	err = store.Set(key, types.STRING, value, time.Hour)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	gotType, gotValue, err := store.GetNoCopy(key)
	if err != nil {
		t.Errorf("GetNoCopy() error = %v", err)
		return
	}

	if gotType != types.STRING {
		t.Errorf("GetNoCopy() gotType = %v, want %v", gotType, types.STRING)
	}

	if string(gotValue) != string(value) {
		t.Errorf("GetNoCopy() gotValue = %v, want %v", string(gotValue), string(value))
	}

	original := string(gotValue)
	gotValue[0] = 'X'

	_, gotValue2, err := store.GetNoCopy(key)
	if err != nil {
		t.Errorf("GetNoCopy() error = %v", err)
		return
	}

	if string(gotValue2) == original {
		t.Error("GetNoCopy() should return reference to original data")
	}
}

func TestCacheStore_Get_EmptyKey(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	_, _, err = store.Get("")
	if err == nil {
		t.Error("Get() should return error for empty key")
	}
	_, _, err = store.GetNoCopy("")
	if err == nil {
		t.Error("GetNoCopy() should return error for empty key")
	}
}

func TestCacheStore_Delete(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	key := "test_key"
	value := []byte("test_value")

	err = store.Set(key, types.STRING, value, time.Hour)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	err = store.Delete(key)
	if err != nil {
		t.Errorf("Delete() error = %v", err)
	}

	_, _, err = store.Get(key)
	if err == nil {
		t.Error("Get() should return error for deleted key")
	}

	err = store.Delete("")
	if err == nil {
		t.Error("Delete() should return error for empty key")
	}
}

func TestCacheStore_Exists(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err = store.Set(key, types.STRING, []byte("value"), time.Hour)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	count := store.Exists(keys...)
	if count != len(keys) {
		t.Errorf("Exists() = %v, want %v", count, len(keys))
	}

	mixedKeys := append(keys, "non_existing")
	count = store.Exists(mixedKeys...)
	if count != len(keys) {
		t.Errorf("Exists() = %v, want %v", count, len(keys))
	}

	count = store.Exists()
	if count != 0 {
		t.Errorf("Exists() = %v, want %v", count, 0)
	}
}

func TestCacheStore_Keys(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	keys := store.Keys()
	if len(keys) != 0 {
		t.Errorf("Keys() = %v, want empty slice", keys)
	}

	expectedKeys := []string{"key1", "key2", "key3"}
	for _, key := range expectedKeys {
		err = store.Set(key, types.STRING, []byte("value"), time.Hour)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	keys = store.Keys()
	if len(keys) != len(expectedKeys) {
		t.Errorf("Keys() length = %v, want %v", len(keys), len(expectedKeys))
	}

	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("Keys() missing key %v", expectedKey)
		}
	}
}

func TestCacheStore_TTL(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	ttl := store.TTL("non_existing")
	if ttl != TTLExpired {
		t.Errorf("TTL() = %v, want %v", ttl, TTLExpired)
	}

	key := "no_expiry_key"
	err = store.Set(key, types.STRING, []byte("value"), 0)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	ttl = store.TTL(key)
	if ttl != TTLNoExpiry {
		t.Errorf("TTL() = %v, want %v", ttl, TTLNoExpiry)
	}

	keyWithExpiry := "expiry_key"
	expiry := time.Hour
	err = store.Set(keyWithExpiry, types.STRING, []byte("value"), expiry)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	ttl = store.TTL(keyWithExpiry)
	if ttl <= 0 || ttl > expiry {
		t.Errorf("TTL() = %v, want between 0 and %v", ttl, expiry)
	}

	expiredKey := "expired_key"
	err = store.Set(expiredKey, types.STRING, []byte("value"), time.Nanosecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(time.Millisecond) // Wait for expiry
	ttl = store.TTL(expiredKey)
	if ttl != TTLExpired {
		t.Errorf("TTL() = %v, want %v", ttl, TTLExpired)
	}
}

func TestCacheStore_Flush(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		err = store.Set(key, types.STRING, []byte("value"), time.Hour)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	count := store.Exists(keys...)
	if count != len(keys) {
		t.Fatalf("Exists() = %v, want %v", count, len(keys))
	}

	store.Flush()

	count = store.Exists(keys...)
	if count != 0 {
		t.Errorf("Exists() after Flush() = %v, want 0", count)
	}
}

func TestCacheStore_CloseAndIsClosed(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	if store.IsClosed() {
		t.Error("IsClosed() should return false for new store")
	}

	err = store.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if !store.IsClosed() {
		t.Error("IsClosed() should return true after Close()")
	}

	err = store.Close()
	if err != nil {
		t.Errorf("Close() second call error = %v", err)
	}
}

func TestCacheStore_ExpiredKeys(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	key := "expired_key"
	err = store.Set(key, types.STRING, []byte("value"), time.Nanosecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(time.Millisecond)

	_, _, err = store.Get(key)
	if err == nil {
		t.Error("Get() should return error for expired key")
	}

	count := store.Exists(key)
	if count != 0 {
		t.Errorf("Exists() = %v, want 0 for expired key", count)
	}
}
