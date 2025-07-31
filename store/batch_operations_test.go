package store

import (
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func TestCacheStore_MGet(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	testData := map[string]struct {
		dataType types.DataType
		value    []byte
		expiry   time.Duration
	}{
		"key1": {types.STRING, []byte("value1"), time.Hour},
		"key2": {types.RAW, []byte("value2"), time.Hour},
		"key3": {types.JSON, []byte(`{"test": "value3"}`), time.Hour},
	}

	for key, data := range testData {
		err := store.Set(key, data.dataType, data.value, data.expiry)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	tests := []struct {
		name     string
		keys     []string
		wantLen  int
		wantErrs []bool
	}{
		{
			name:     "get existing keys",
			keys:     []string{"key1", "key2", "key3"},
			wantLen:  3,
			wantErrs: []bool{false, false, false},
		},
		{
			name:     "get mix of existing and non-existing keys",
			keys:     []string{"key1", "nonexistent", "key2"},
			wantLen:  3,
			wantErrs: []bool{false, true, false},
		},
		{
			name:     "get non-existing keys",
			keys:     []string{"nonexistent1", "nonexistent2"},
			wantLen:  2,
			wantErrs: []bool{true, true},
		},
		{
			name:     "empty key",
			keys:     []string{""},
			wantLen:  1,
			wantErrs: []bool{true},
		},
		{
			name:    "no keys",
			keys:    []string{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := store.MGet(tt.keys...)

			if tt.wantLen == 0 && results == nil {
				return
			}

			if len(results) != tt.wantLen {
				t.Errorf("MGet() returned %d results, want %d", len(results), tt.wantLen)
				return
			}

			for i, result := range results {
				if result.Key != tt.keys[i] {
					t.Errorf("MGet() result[%d].Key = %v, want %v", i, result.Key, tt.keys[i])
				}

				hasError := result.Error != nil
				wantError := tt.wantErrs[i]

				if hasError != wantError {
					t.Errorf("MGet() result[%d] error = %v, wantError %v", i, result.Error, wantError)
				}

				if !hasError {
					expectedData := testData[tt.keys[i]]
					if result.Type != expectedData.dataType {
						t.Errorf("MGet() result[%d].Type = %v, want %v", i, result.Type, expectedData.dataType)
					}
					if string(result.Value) != string(expectedData.value) {
						t.Errorf("MGet() result[%d].Value = %v, want %v", i, string(result.Value), string(expectedData.value))
					}
				}
			}
		})
	}
}

func TestCacheStore_MGet_ExpiredKeys(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	err = store.Set("expired_key", types.STRING, []byte("value"), time.Nanosecond)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	time.Sleep(time.Millisecond)

	results := store.MGet("expired_key")
	if len(results) != 1 {
		t.Fatalf("MGet() returned %d results, want 1", len(results))
	}

	if results[0].Error == nil {
		t.Error("MGet() should return error for expired key")
	}
}

func TestCacheStore_MSet(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name      string
		items     []BatchItem
		wantErrs  []bool
		wantCount int
	}{
		{
			name: "valid items",
			items: []BatchItem{
				NewItem("key1", types.STRING, []byte("value1"), time.Hour),
				NewItem("key2", types.RAW, []byte("value2"), time.Hour),
				NewItem("key3", types.JSON, []byte(`{"test": "value3"}`), time.Hour),
			},
			wantErrs:  []bool{false, false, false},
			wantCount: 0,
		},
		{
			name: "mixed valid and invalid items",
			items: []BatchItem{
				NewItem("key1", types.STRING, []byte("value1"), time.Hour),
				NewItem("", types.STRING, []byte("value2"), time.Hour),
				NewItem("key3", types.STRING, nil, time.Hour),
				NewItem("key4", types.RAW, []byte("value4"), time.Hour),
			},
			wantErrs:  []bool{false, true, true, false},
			wantCount: 2,
		},
		{
			name: "empty key",
			items: []BatchItem{
				NewItem("", types.STRING, []byte("value"), time.Hour),
			},
			wantErrs:  []bool{true},
			wantCount: 1,
		},
		{
			name: "nil value",
			items: []BatchItem{
				NewItem("key", types.STRING, nil, time.Hour),
			},
			wantErrs:  []bool{true},
			wantCount: 1,
		},
		{
			name:      "no items",
			items:     []BatchItem{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := store.MSet(tt.items...)

			if tt.wantCount == 0 && errs == nil {
				return
			}

			if len(errs) != len(tt.items) {
				t.Errorf("MSet() returned %d errors, want %d", len(errs), len(tt.items))
				return
			}

			errorCount := 0
			for i, err := range errs {
				hasError := err != nil
				wantError := tt.wantErrs[i]

				if hasError != wantError {
					t.Errorf("MSet() error[%d] = %v, wantError %v", i, err, wantError)
				}

				if hasError {
					errorCount++
				}
			}

			if errorCount != tt.wantCount {
				t.Errorf("MSet() error count = %d, want %d", errorCount, tt.wantCount)
			}

			for i, item := range tt.items {
				if errs[i] == nil {
					dataType, value, err := store.Get(item.Key)
					if err != nil {
						t.Errorf("Get() after MSet() error = %v  %v", err, store.memorydbTemporary)
						continue
					}
					if dataType != item.Entry.Type {
						t.Errorf("Get() after MSet() type = %v, want %v", dataType, item.Entry.Type)
					}
					if string(value) != string(item.Entry.Data) {
						t.Errorf("Get() after MSet() value = %v, want %v", string(value), string(item.Entry.Data))
					}
				}
			}
		})
	}
}

func TestCacheStore_MDelete(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	testKeys := []string{"key1", "key2", "key3"}
	for _, key := range testKeys {
		err := store.Set(key, types.STRING, []byte("value"), time.Hour)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	tests := []struct {
		name     string
		keys     []string
		wantErrs []bool
	}{
		{
			name:     "delete existing keys",
			keys:     []string{"key1", "key2"},
			wantErrs: []bool{false, false},
		},
		{
			name:     "delete mix of existing and non-existing keys",
			keys:     []string{"key3", "nonexistent"},
			wantErrs: []bool{false, false},
		},
		{
			name:     "delete empty key",
			keys:     []string{""},
			wantErrs: []bool{true},
		},
		{
			name: "no keys",
			keys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := store.MDelete(tt.keys...)

			if len(tt.keys) == 0 && errs == nil {
				return
			}

			if len(errs) != len(tt.keys) {
				t.Errorf("MDelete() returned %d errors, want %d", len(errs), len(tt.keys))
				return
			}

			for i, err := range errs {
				hasError := err != nil
				wantError := tt.wantErrs[i]

				if hasError != wantError {
					t.Errorf("MDelete() error[%d] = %v, wantError %v", i, err, wantError)
				}

				if tt.keys[i] != "" && err == nil {
					_, _, getErr := store.Get(tt.keys[i])
					if getErr == nil {
						t.Errorf("Get() after MDelete() should return error for deleted key %v", tt.keys[i])
					}
				}
			}
		})
	}
}

func TestCacheStore_MDelete_EmptyKey(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	errs := store.MDelete("")
	if len(errs) != 1 {
		t.Fatalf("MDelete() returned %d errors, want 1", len(errs))
	}

	if errs[0] != errors.ErrKeyEmpty {
		t.Errorf("MDelete() error = %v, want %v", errs[0], errors.ErrKeyEmpty)
	}
}

func TestBatchOperations_Integration(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	items := []BatchItem{
		NewItem("user:1", types.JSON, []byte(`{"name": "Alice", "age": 30}`), time.Hour),
		NewItem("user:2", types.JSON, []byte(`{"name": "Bob", "age": 25}`), time.Hour),
		NewItem("counter", types.RAW, []byte("100"), time.Hour),
		NewItem("flag", types.STRING, []byte("enabled"), time.Hour),
	}

	errs := store.MSet(items...)
	for i, err := range errs {
		if err != nil {
			t.Errorf("MSet() error[%d] = %v", i, err)
		}
	}

	keys := []string{"user:1", "user:2", "counter", "flag", "nonexistent"}
	results := store.MGet(keys...)

	if len(results) != 5 {
		t.Fatalf("MGet() returned %d results, want 5", len(results))
	}

	for i := 0; i < 4; i++ {
		if results[i].Error != nil {
			t.Errorf("MGet() result[%d] error = %v", i, results[i].Error)
		}
		if string(results[i].Value) != string(items[i].Entry.Data) {
			t.Errorf("MGet() result[%d] value = %v, want %v", i, string(results[i].Value), string(items[i].Entry.Data))
		}
	}

	if results[4].Error == nil {
		t.Error("MGet() should return error for nonexistent key")
	}

	deleteKeys := []string{"user:1", "counter"}
	deleteErrs := store.MDelete(deleteKeys...)
	for i, err := range deleteErrs {
		if err != nil {
			t.Errorf("MDelete() error[%d] = %v", i, err)
		}
	}

	afterDeleteResults := store.MGet(deleteKeys...)
	for i, result := range afterDeleteResults {
		if result.Error == nil {
			t.Errorf("MGet() after delete should return error for key %v", deleteKeys[i])
		}
	}

	remainingResults := store.MGet("user:2", "flag")
	for i, result := range remainingResults {
		if result.Error != nil {
			t.Errorf("MGet() remaining key[%d] error = %v", i, result.Error)
		}
	}
}
