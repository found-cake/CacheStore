package store

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

func TestCacheStore_BoolOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name    string
		key     string
		value   bool
		expiry  time.Duration
		wantErr bool
	}{
		{"set true", "bool_true", true, time.Hour, false},
		{"set false", "bool_false", false, time.Hour, false},
		{"empty key", "", true, time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetBool(tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := store.GetBool(tt.key)
				if err != nil {
					t.Errorf("GetBool() error = %v", err)
					return
				}
				if got != tt.value {
					t.Errorf("GetBool() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestCacheStore_IntegerOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("Int16", func(t *testing.T) {
		key := "int16_key"
		value := int16(12345)

		err := store.SetInt16(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetInt16() error = %v", err)
		}

		got, err := store.GetInt16(key)
		if err != nil {
			t.Errorf("GetInt16() error = %v", err)
		}
		if got != value {
			t.Errorf("GetInt16() = %v, want %v", got, value)
		}

		delta := int16(100)
		err = store.IncrInt16(key, delta, 0)
		if err != nil {
			t.Errorf("IncrInt16() error = %v", err)
		}

		got, err = store.GetInt16(key)
		if err != nil {
			t.Errorf("GetInt16() after incr error = %v", err)
		}
		expected := value + delta
		if got != expected {
			t.Errorf("GetInt16() after incr = %v, want %v", got, expected)
		}
	})

	t.Run("Int32", func(t *testing.T) {
		key := "int32_key"
		value := int32(1234567890)

		err := store.SetInt32(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetInt32() error = %v", err)
		}

		got, err := store.GetInt32(key)
		if err != nil {
			t.Errorf("GetInt32() error = %v", err)
		}
		if got != value {
			t.Errorf("GetInt32() = %v, want %v", got, value)
		}

		delta := int32(1000)
		err = store.IncrInt32(key, delta, 0)
		if err != nil {
			t.Errorf("IncrInt32() error = %v", err)
		}

		got, err = store.GetInt32(key)
		if err != nil {
			t.Errorf("GetInt32() after incr error = %v", err)
		}
		expected := value + delta
		if got != expected {
			t.Errorf("GetInt32() after incr = %v, want %v", got, expected)
		}
	})

	t.Run("Int64", func(t *testing.T) {
		key := "int64_key"
		value := int64(1234567890123456789)

		err := store.SetInt64(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetInt64() error = %v", err)
		}

		got, err := store.GetInt64(key)
		if err != nil {
			t.Errorf("GetInt64() error = %v", err)
		}
		if got != value {
			t.Errorf("GetInt64() = %v, want %v", got, value)
		}

		delta := int64(10000)
		err = store.IncrInt64(key, delta, 0)
		if err != nil {
			t.Errorf("IncrInt64() error = %v", err)
		}

		got, err = store.GetInt64(key)
		if err != nil {
			t.Errorf("GetInt64() after incr error = %v", err)
		}
		expected := value + delta
		if got != expected {
			t.Errorf("GetInt64() after incr = %v, want %v", got, expected)
		}
	})
}

func TestCacheStore_UnsignedIntegerOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("UInt16", func(t *testing.T) {
		key := "uint16_key"
		value := uint16(65535)

		err := store.SetUInt16(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetUInt16() error = %v", err)
		}

		got, err := store.GetUInt16(key)
		if err != nil {
			t.Errorf("GetUInt16() error = %v", err)
		}
		if got != value {
			t.Errorf("GetUInt16() = %v, want %v", got, value)
		}

		delta := uint16(100)
		err = store.IncrUInt16(key, delta, 0)
		if err == nil {
			t.Error("IncrUInt16() should return overflow error")
		}

		err = store.SetUInt16(key, 1000, time.Hour)
		if err != nil {
			t.Errorf("SetUInt16() error = %v", err)
		}

		err = store.DecrUInt16(key, 500, 0)
		if err != nil {
			t.Errorf("DecrUInt16() error = %v", err)
		}

		got, err = store.GetUInt16(key)
		if err != nil {
			t.Errorf("GetUInt16() after decr error = %v", err)
		}
		if got != 500 {
			t.Errorf("GetUInt16() after decr = %v, want 500", got)
		}

		err = store.DecrUInt16(key, 1000, 0)
		if err == nil {
			t.Error("DecrUInt16() should return underflow error")
		}
	})

	t.Run("UInt32", func(t *testing.T) {
		key := "uint32_key"
		value := uint32(4294967295)

		err := store.SetUInt32(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetUInt32() error = %v", err)
		}

		got, err := store.GetUInt32(key)
		if err != nil {
			t.Errorf("GetUInt32() error = %v", err)
		}
		if got != value {
			t.Errorf("GetUInt32() = %v, want %v", got, value)
		}

		err = store.SetUInt32(key, 100, time.Hour)
		if err != nil {
			t.Errorf("SetUInt32() error = %v", err)
		}

		delta := uint32(50)
		err = store.IncrUInt32(key, delta, 0)
		if err != nil {
			t.Errorf("IncrUInt32() error = %v", err)
		}

		got, err = store.GetUInt32(key)
		if err != nil {
			t.Errorf("GetUInt32() after incr error = %v", err)
		}
		if got != 150 {
			t.Errorf("GetUInt32() after incr = %v, want 150", got)
		}
	})

	t.Run("UInt64", func(t *testing.T) {
		key := "uint64_key"
		value := uint64(18446744073709551615)

		err := store.SetUInt64(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetUInt64() error = %v", err)
		}

		got, err := store.GetUInt64(key)
		if err != nil {
			t.Errorf("GetUInt64() error = %v", err)
		}
		if got != value {
			t.Errorf("GetUInt64() = %v, want %v", got, value)
		}

		err = store.SetUInt64(key, 1000, time.Hour)
		if err != nil {
			t.Errorf("SetUInt64() error = %v", err)
		}

		err = store.IncrUInt64(key, 500, 0)
		if err != nil {
			t.Errorf("IncrUInt64() error = %v", err)
		}

		got, err = store.GetUInt64(key)
		if err != nil {
			t.Errorf("GetUInt64() after incr error = %v", err)
		}
		if got != 1500 {
			t.Errorf("GetUInt64() after incr = %v, want 1500", got)
		}

		err = store.DecrUInt64(key, 300, 0)
		if err != nil {
			t.Errorf("DecrUInt64() error = %v", err)
		}

		got, err = store.GetUInt64(key)
		if err != nil {
			t.Errorf("GetUInt64() after decr error = %v", err)
		}
		if got != 1200 {
			t.Errorf("GetUInt64() after decr = %v, want 1200", got)
		}
	})
}

func TestCacheStore_FloatOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	t.Run("Float32", func(t *testing.T) {
		key := "float32_key"
		value := float32(3.14159)

		err := store.SetFloat32(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetFloat32() error = %v", err)
		}

		got, err := store.GetFloat32(key)
		if err != nil {
			t.Errorf("GetFloat32() error = %v", err)
		}
		if math.Abs(float64(got-value)) > 1e-6 {
			t.Errorf("GetFloat32() = %v, want %v", got, value)
		}

		delta := float32(1.5)
		err = store.IncrFloat32(key, delta, 0)
		if err != nil {
			t.Errorf("IncrFloat32() error = %v", err)
		}

		got, err = store.GetFloat32(key)
		if err != nil {
			t.Errorf("GetFloat32() after incr error = %v", err)
		}
		expected := value + delta
		if math.Abs(float64(got-expected)) > 1e-6 {
			t.Errorf("GetFloat32() after incr = %v, want %v", got, expected)
		}
	})

	t.Run("Float64", func(t *testing.T) {
		key := "float64_key"
		value := float64(3.141592653589793)

		err := store.SetFloat64(key, value, time.Hour)
		if err != nil {
			t.Errorf("SetFloat64() error = %v", err)
		}

		got, err := store.GetFloat64(key)
		if err != nil {
			t.Errorf("GetFloat64() error = %v", err)
		}
		if math.Abs(got-value) > 1e-15 {
			t.Errorf("GetFloat64() = %v, want %v", got, value)
		}

		delta := float64(2.5)
		err = store.IncrFloat64(key, delta, 0)
		if err != nil {
			t.Errorf("IncrFloat64() error = %v", err)
		}

		got, err = store.GetFloat64(key)
		if err != nil {
			t.Errorf("GetFloat64() after incr error = %v", err)
		}
		expected := value + delta
		if math.Abs(got-expected) > 1e-15 {
			t.Errorf("GetFloat64() after incr = %v, want %v", got, expected)
		}
	})

	t.Run("SpecialValues", func(t *testing.T) {
		tests := []struct {
			name  string
			value float64
		}{
			{"positive infinity", math.Inf(1)},
			{"negative infinity", math.Inf(-1)},
			{"NaN", math.NaN()},
			{"positive zero", 0.0},
			{"negative zero", math.Copysign(0.0, -1)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				key := "special_" + tt.name
				err := store.SetFloat64(key, tt.value, time.Hour)
				if err != nil {
					t.Errorf("SetFloat64() error = %v", err)
				}

				got, err := store.GetFloat64(key)
				if err != nil {
					t.Errorf("GetFloat64() error = %v", err)
				}

				if math.IsNaN(tt.value) {
					if !math.IsNaN(got) {
						t.Errorf("GetFloat64() = %v, want NaN", got)
					}
				} else if got != tt.value {
					t.Errorf("GetFloat64() = %v, want %v", got, tt.value)
				}
			})
		}
	})
}

func TestCacheStore_StringOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name    string
		key     string
		value   string
		expiry  time.Duration
		wantErr bool
	}{
		{"normal string", "str1", "hello world", time.Hour, false},
		{"empty string", "str2", "", time.Hour, false},
		{"unicode string", "str3", "こんにちは 世界", time.Hour, false},
		{"empty key", "", "value", time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetString(tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := store.GetString(tt.key)
				if err != nil {
					t.Errorf("GetString() error = %v", err)
					return
				}
				if got != tt.value {
					t.Errorf("GetString() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestCacheStore_RawOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name    string
		key     string
		value   []byte
		expiry  time.Duration
		wantErr bool
	}{
		{"normal bytes", "raw1", []byte("binary data"), time.Hour, false},
		{"empty bytes", "raw2", []byte{}, time.Hour, false},
		{"binary data", "raw3", []byte{0x00, 0x01, 0xFF, 0xFE}, time.Hour, false},
		{"empty key", "", []byte("data"), time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetRaw(tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := store.GetRaw(tt.key)
				if err != nil {
					t.Errorf("GetRaw() error = %v", err)
					return
				}
				if string(got) != string(tt.value) {
					t.Errorf("GetRaw() = %v, want %v", got, tt.value)
				}

				gotNoCopy, err := store.GetRawNoCopy(tt.key)
				if err != nil {
					t.Errorf("GetRawNoCopy() error = %v", err)
					return
				}
				if string(gotNoCopy) != string(tt.value) {
					t.Errorf("GetRawNoCopy() = %v, want %v", gotNoCopy, tt.value)
				}

				if len(gotNoCopy) > 0 {
					original := gotNoCopy[0]
					gotNoCopy[0] = 'X'

					gotAgain, err := store.GetRawNoCopy(tt.key)
					if err != nil {
						t.Errorf("GetRawNoCopy() second call error = %v", err)
						return
					}
					if len(gotAgain) > 0 && gotAgain[0] == original {
						t.Error("GetRawNoCopy() should return reference to original data")
					}
				}
			}
		})
	}
}

func TestCacheStore_JSONOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	type TestStruct struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	tests := []struct {
		name    string
		key     string
		value   interface{}
		expiry  time.Duration
		wantErr bool
	}{
		{
			"struct",
			"json1",
			TestStruct{Name: "Alice", Age: 30, Email: "alice@example.com"},
			time.Hour,
			false,
		},
		{
			"map",
			"json2",
			map[string]interface{}{"key": "value", "number": 42},
			time.Hour,
			false,
		},
		{
			"slice",
			"json3",
			[]string{"a", "b", "c"},
			time.Hour,
			false,
		},
		{
			"empty key",
			"",
			TestStruct{},
			time.Hour,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetJSON(tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				switch v := tt.value.(type) {
				case TestStruct:
					var got TestStruct
					err := store.GetJSON(tt.key, &got)
					if err != nil {
						t.Errorf("GetJSON() error = %v", err)
						return
					}
					if got != v {
						t.Errorf("GetJSON() = %v, want %v", got, v)
					}

				case map[string]interface{}:
					var got map[string]interface{}
					err := store.GetJSON(tt.key, &got)
					if err != nil {
						t.Errorf("GetJSON() error = %v", err)
						return
					}
					expectedJSON, _ := json.Marshal(v)
					gotJSON, _ := json.Marshal(got)
					if string(gotJSON) != string(expectedJSON) {
						t.Errorf("GetJSON() = %v, want %v", got, v)
					}

				case []string:
					var got []string
					err := store.GetJSON(tt.key, &got)
					if err != nil {
						t.Errorf("GetJSON() error = %v", err)
						return
					}
					if len(got) != len(v) {
						t.Errorf("GetJSON() length = %d, want %d", len(got), len(v))
						return
					}
					for i, item := range got {
						if item != v[i] {
							t.Errorf("GetJSON()[%d] = %v, want %v", i, item, v[i])
						}
					}
				}
			}
		})
	}

	err = store.Set("invalid_json", types.JSON, []byte("{invalid json"), time.Hour)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	var target TestStruct
	err = store.GetJSON("invalid_json", &target)
	if err == nil {
		t.Error("GetJSON() should return error for invalid JSON")
	}

	type BadTarget struct {
		Ch chan int `json:"channel"`
	}
	badTarget := BadTarget{Ch: make(chan int)}
	err = store.SetJSON("bad_json", badTarget, time.Hour)
	if err == nil {
		t.Error("SetJSON() should return error for unmarshalable data")
	}
}

func TestCacheStore_TimeOperations(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()
	utc := now.UTC()
	local := now.Local()

	tests := []struct {
		name    string
		key     string
		value   time.Time
		expiry  time.Duration
		wantErr bool
	}{
		{"current time", "time1", now, time.Hour, false},
		{"utc time", "time2", utc, time.Hour, false},
		{"local time", "time3", local, time.Hour, false},
		{"zero time", "time4", time.Time{}, time.Hour, false},
		{"empty key", "", now, time.Hour, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SetTime(tt.key, tt.value, tt.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got, err := store.GetTime(tt.key)
				if err != nil {
					t.Errorf("GetTime() error = %v", err)
					return
				}
				if !got.Equal(tt.value) {
					t.Errorf("GetTime() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestCacheStore_TypeMismatchErrors(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	err = store.SetString("string_key", "hello", time.Hour)
	if err != nil {
		t.Fatalf("SetString() error = %v", err)
	}

	err = store.SetInt32("int_key", 42, time.Hour)
	if err != nil {
		t.Fatalf("SetInt32() error = %v", err)
	}

	err = store.SetBool("bool_key", true, time.Hour)
	if err != nil {
		t.Fatalf("SetBool() error = %v", err)
	}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
	}{
		{
			"string as int",
			func() error {
				_, err := store.GetInt32("string_key")
				return err
			},
			true,
		},
		{
			"int as string",
			func() error {
				_, err := store.GetString("int_key")
				return err
			},
			true,
		},
		{
			"bool as int",
			func() error {
				_, err := store.GetInt32("bool_key")
				return err
			},
			true,
		},
		{
			"string as bool",
			func() error {
				_, err := store.GetBool("string_key")
				return err
			},
			true,
		},
		{
			"int as float",
			func() error {
				_, err := store.GetFloat32("int_key")
				return err
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if (err != nil) != tt.expectError {
				t.Errorf("Operation error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestCacheStore_EmptyKeyErrors(t *testing.T) {
	store, err := NewCacheStore(config.Config{DBSave: false})
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	tests := []struct {
		name      string
		operation func() error
	}{
		{"GetString", func() error { _, err := store.GetString(""); return err }},
		{"GetBool", func() error { _, err := store.GetBool(""); return err }},
		{"GetRaw", func() error { _, err := store.GetRaw(""); return err }},
		{"GetRawNoCopy", func() error { _, err := store.GetRawNoCopy(""); return err }},
		{"GetJSON", func() error { var v interface{}; return store.GetJSON("", &v) }},
		{"GetTime", func() error { _, err := store.GetTime(""); return err }},
		{"GetInt16", func() error { _, err := store.GetInt16(""); return err }},
		{"GetInt32", func() error { _, err := store.GetInt32(""); return err }},
		{"GetInt64", func() error { _, err := store.GetInt64(""); return err }},
		{"GetUInt16", func() error { _, err := store.GetUInt16(""); return err }},
		{"GetUInt32", func() error { _, err := store.GetUInt32(""); return err }},
		{"GetUInt64", func() error { _, err := store.GetUInt64(""); return err }},
		{"GetFloat32", func() error { _, err := store.GetFloat32(""); return err }},
		{"GetFloat64", func() error { _, err := store.GetFloat64(""); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()
			if err != errors.ErrKeyEmpty {
				t.Errorf("%s with empty key error = %v, want %v", tt.name, err, errors.ErrKeyEmpty)
			}
		})
	}
}
