package store

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/utils/types"
)

func BenchmarkCurrentStructure_GoroutineConcurrency(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 500 * time.Millisecond,
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

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 100
		operationsPerGoroutine := 100

		wg.Add(numGoroutines)

		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < operationsPerGoroutine; k++ {
					key := keys[rand.Intn(len(keys))]

					if rand.Float32() < 0.7 {
						_, _, _ = store.Get(key)
					} else {
						value := make([]byte, 100)
						cr.Read(value)
						_ = store.Set(key, types.RAW, value, time.Hour)
					}
				}
			}()
		}

		wg.Wait()
	}
}

func BenchmarkCurrentStructure_LockContentionGoroutines(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 100 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	permanentKeys := make([]string, 50)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("permanent_%d", i)
		permanentKeys[i] = key
		value := make([]byte, 100)
		cr.Read(value)
		store.Set(key, types.RAW, value, 0)
	}

	for i := 0; i < 500; i++ {
		key := fmt.Sprintf("temp_%d", i)
		value := make([]byte, 100)
		cr.Read(value)
		store.Set(key, types.RAW, value, 50*time.Millisecond)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numReaders := 50
		numWriters := 10

		wg.Add(numReaders)
		for j := 0; j < numReaders; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < 20; k++ {
					key := permanentKeys[rand.Intn(len(permanentKeys))]
					_, _, _ = store.Get(key)
					time.Sleep(time.Microsecond)
				}
			}()
		}

		wg.Add(numWriters)
		for j := 0; j < numWriters; j++ {
			go func(goroutineID int) {
				defer wg.Done()
				for k := 0; k < 10; k++ {
					key := fmt.Sprintf("new_temp_%d_%d", goroutineID, k)
					value := make([]byte, 100)
					cr.Read(value)
					_ = store.Set(key, types.RAW, value, 100*time.Millisecond)
					time.Sleep(time.Microsecond)
				}
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkCurrentStructure_GCDuringConcurrentAccess(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 50 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	permanentKeys := make([]string, 100)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("permanent_%d", i)
		permanentKeys[i] = key
		value := make([]byte, 200)
		cr.Read(value)
		store.Set(key, types.RAW, value, 0)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("temp_gc_%d", j)
				value := make([]byte, 100)
				cr.Read(value)
				store.Set(key, types.RAW, value, 10*time.Millisecond)
				time.Sleep(time.Millisecond)
			}
		}()

		numReaders := 20
		wg.Add(numReaders)
		for j := 0; j < numReaders; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < 50; k++ {
					key := permanentKeys[rand.Intn(len(permanentKeys))]
					_, _, _ = store.Get(key)
				}
			}()
		}

		wg.Wait()
	}
}

func BenchmarkCurrentStructure_ReadWriteRatios(b *testing.B) {
	ratios := []struct {
		name      string
		readRatio float32
	}{
		{"Read100_Write0", 1.0},
		{"Read90_Write10", 0.9},
		{"Read80_Write20", 0.8},
		{"Read70_Write30", 0.7},
		{"Read50_Write50", 0.5},
	}

	for _, ratio := range ratios {
		b.Run(ratio.name, func(b *testing.B) {
			store, err := NewCacheStore(config.Config{
				DBSave:     false,
				GCInterval: 1 * time.Second,
			})
			if err != nil {
				b.Fatal(err)
			}
			defer store.Close()

			keys := make([]string, 500)
			for i := 0; i < 500; i++ {
				key := fmt.Sprintf("key_%d", i)
				keys[i] = key
				value := make([]byte, 100)
				cr.Read(value)
				store.Set(key, types.RAW, value, time.Hour)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				var wg sync.WaitGroup
				numGoroutines := 20
				operationsPerGoroutine := 50

				wg.Add(numGoroutines)
				for j := 0; j < numGoroutines; j++ {
					go func() {
						defer wg.Done()
						for k := 0; k < operationsPerGoroutine; k++ {
							key := keys[rand.Intn(len(keys))]

							if rand.Float32() < ratio.readRatio {
								_, _, _ = store.Get(key)
							} else {
								value := make([]byte, 100)
								cr.Read(value)
								_ = store.Set(key, types.RAW, value, time.Hour)
							}
						}
					}()
				}
				wg.Wait()
			}
		})
	}
}

func BenchmarkCurrentStructure_MassiveGoroutineStress(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 200 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	baseKeys := make([]string, 200)
	for i := 0; i < 200; i++ {
		key := fmt.Sprintf("base_%d", i)
		baseKeys[i] = key
		value := make([]byte, 150)
		cr.Read(value)

		if i%2 == 0 {
			store.Set(key, types.RAW, value, 0)
		} else {
			store.Set(key, types.RAW, value, 5*time.Minute)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 200

		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func(goroutineID int) {
				defer wg.Done()

				for k := 0; k < 10; k++ {
					key := baseKeys[rand.Intn(len(baseKeys))]
					_, _, _ = store.Get(key)

					if k%5 == 0 {
						newKey := fmt.Sprintf("temp_%d_%d", goroutineID, k)
						value := make([]byte, 100)
						cr.Read(value)
						_ = store.Set(newKey, types.RAW, value, 1*time.Minute)
					}

					if k%7 == 0 {
						deleteKey := fmt.Sprintf("temp_%d_%d", goroutineID, k-1)
						_ = store.Delete(deleteKey)
					}
				}
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkCurrentStructure_TTLConcurrency(b *testing.B) {
	store, err := NewCacheStore(config.Config{
		DBSave:     false,
		GCInterval: 300 * time.Millisecond,
	})
	if err != nil {
		b.Fatal(err)
	}
	defer store.Close()

	keys := make([]string, 300)
	for i := 0; i < 300; i++ {
		key := fmt.Sprintf("ttl_key_%d", i)
		keys[i] = key
		value := make([]byte, 100)
		cr.Read(value)

		var expiry time.Duration
		switch i % 4 {
		case 0:
			expiry = 0
		case 1:
			expiry = 1 * time.Second
		case 2:
			expiry = 5 * time.Second
		case 3:
			expiry = 10 * time.Second
		}

		store.Set(key, types.RAW, value, expiry)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 30

		wg.Add(numGoroutines)
		for j := 0; j < numGoroutines; j++ {
			go func() {
				defer wg.Done()
				for k := 0; k < 20; k++ {
					key := keys[rand.Intn(len(keys))]
					_ = store.TTL(key)
					_, _, _ = store.Get(key)
					_ = store.Exists(key)
				}
			}()
		}

		wg.Wait()
	}
}
