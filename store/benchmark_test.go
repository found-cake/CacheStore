package store

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func generateTestData(count int) map[string][]byte {
	data := make(map[string][]byte, count)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := make([]byte, 100)
		cr.Read(value)
		data[key] = value
	}
	return data
}

func BenchmarkCurrentStructure_Get(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	testData := generateTestData(1000)
	keys := make([]string, 0, len(testData))

	i := 0
	for key, value := range testData {
		if i%2 == 0 {
			store.Set(key, types.RAW, value, 0)
		} else {
			store.Set(key, types.RAW, value, time.Hour)
		}
		keys = append(keys, key)
		i++
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(len(keys))]
			_, _, _ = store.Get(key)
		}
	})
}

func BenchmarkCurrentStructure_GetDuringGC(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 100 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	testData := generateTestData(2000)
	keys := make([]string, 0, len(testData))

	for key, value := range testData {
		store.Set(key, types.RAW, value, 50*time.Millisecond)
		keys = append(keys, key)
	}

	permanentData := generateTestData(100)
	for key, value := range permanentData {
		store.Set(key, types.RAW, value, 0)
		keys = append(keys, key)
	}

	time.Sleep(200 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(len(keys))]
			_, _, _ = store.Get(key)
		}
	})
}

func BenchmarkCurrentStructure_ConcurrentAccess(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	testData := generateTestData(500)
	keys := make([]string, 0, len(testData))

	for key, value := range testData {
		store.Set(key, types.RAW, value, time.Hour)
		keys = append(keys, key)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Float32() < 0.8 {
				key := keys[rand.Intn(len(keys))]
				_, _, _ = store.Get(key)
			} else {
				key := keys[rand.Intn(len(keys))]
				value := make([]byte, 100)
				cr.Read(value)
				_ = store.Set(key, types.RAW, value, time.Hour)
			}
		}
	})
}

func BenchmarkCurrentStructure_LockContention(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 50 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("temp_key_%d", i)
		value := make([]byte, 50)
		cr.Read(value)
		store.Set(key, types.RAW, value, 100*time.Millisecond)
	}

	permanentKeys := make([]string, 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("permanent_key_%d", i)
		permanentKeys[i] = key
		value := make([]byte, 50)
		cr.Read(value)
		store.Set(key, types.RAW, value, 0)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := permanentKeys[rand.Intn(len(permanentKeys))]
			_, _, _ = store.Get(key)
		}
	})
}

func BenchmarkCurrentStructure_MemoryProfile(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := make([]byte, 1024)
		cr.Read(value)

		if i%3 == 0 {
			store.Set(key, types.RAW, value, 0)
		} else {
			store.Set(key, types.RAW, value, time.Hour)
		}
	}
}

func BenchmarkCurrentStructure_RealWorldScenario(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 500 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	configKeys := make([]string, 50)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("config_%d", i)
		configKeys[i] = key
		value := []byte(fmt.Sprintf("config_value_%d", i))
		store.Set(key, types.STRING, value, 0)
	}

	sessionKeys := make([]string, 200)
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("session_%d", i)
		sessionKeys[i] = key
		value := make([]byte, 256)
		cr.Read(value)
		store.Set(key, types.RAW, value, 30*time.Minute)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.Float32() < 0.7 {
				key := configKeys[rand.Intn(len(configKeys))]
				_, _, _ = store.Get(key)
			} else {
				key := sessionKeys[rand.Intn(len(sessionKeys))]
				_, _, _ = store.Get(key)
			}
		}
	})
}

func BenchmarkCurrentStructure_GetNoCopy(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 1 * time.Second,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	testData := generateTestData(1000)
	keys := make([]string, 0, len(testData))

	for key, value := range testData {
		store.Set(key, types.RAW, value, time.Hour)
		keys = append(keys, key)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(len(keys))]
			_, _, _ = store.GetNoCopy(key)
		}
	})
}

func BenchmarkCurrentStructure_CleanExpired(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 0,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("expired_key_%d", i)
		value := make([]byte, 100)
		cr.Read(value)
		store.Set(key, types.RAW, value, 1*time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.cleanExpired()
	}
}
