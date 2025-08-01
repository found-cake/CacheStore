# CacheStore

[![Go Test](https://github.com/found-cake/CacheStore/actions/workflows/gotest.yml/badge.svg)](https://github.com/found-cake/CacheStore/actions/workflows/gotest.yml)
[![Go Version](https://img.shields.io/badge/go-1.22%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

> CacheStore is a thread-safe cache library written in Go, providing blazing-fast memory access and reliable SQLite-based data persistence.

## ✨ Features

- 💾 **SQLite persistence**: Durable data storage and recovery
- 🔒 **Thread-safe**: Concurrency safety with RWMutex
- ⏰ **TTL support**: Automatic expiration and garbage collection
- 📊 **Various data types**: String, JSON, Boolean, Integer (16/32/64bit), Time
- 🚀 **Batch operations**: Supports MGet, MSet, MDelete
- 🎯 **Dirty data management**: Smart change tracking and sync
- ⚡ **Zero-copy option**: Performance-optimized GetNoCopy method

## 📦 Installation

```bash
go get github.com/found-cake/CacheStore
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/found-cake/CacheStore/config"
    "github.com/found-cake/CacheStore/store"
    "github.com/found-cake/CacheStore/utils/types"
)

func main() {
    // Create a cache store with default config
    cfg := config.DefaultConfig()
    cacheStore, err := store.NewCacheStore(cfg)
    if err != nil {
        panic(err)
    }
    defer cacheStore.Close()

    // Set string value
    err = cacheStore.Set("user:123", types.STRING, []byte("Alice"), time.Hour)
    if err != nil {
        panic(err)
    }

    // Get value
    dataType, value, err := cacheStore.Get("user:123")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Type: %v, Value: %s\n", dataType, string(value))
}
```

### Batch Operations

```go
import "github.com/found-cake/CacheStore/store"

// Batch set
items := []store.BatchItem{
    store.NewItem("key1", types.STRING, []byte("value1"), time.Hour),
    store.NewItem("key2", types.STRING, []byte("value2"), time.Hour),
    store.NewItem("key3", types.STRING, []byte("value3"), time.Hour),
}

errs := cacheStore.MSet(items...)
for i, err := range errs {
    if err != nil {
        fmt.Printf("Error setting item %d: %v\n", i, err)
    }
}

// Batch get
keys := []string{"key1", "key2", "key3"}
results := cacheStore.MGet(keys...)
for i, result := range results {
    if result.Error != nil {
        fmt.Printf("Error getting %s: %v\n", keys[i], result.Error)
    } else {
        fmt.Printf("%s = %s\n", keys[i], string(result.Value))
    }
}
```

## ⚙️ Configuration Options

```go
cfg := config.Config{
    GCInterval:          10 * time.Second,  // Garbage collection interval
    DBSave:              true,              // Enable SQLite persistence
    DBFileName:          "cache.db",        // Database file name
    DBSaveInterval:      10 * time.Minute,  // DB save interval
    SaveDirtyData:       true,              // Enable dirty data tracking
    DirtyThresholdCount: 50,                // Dirty item threshold
    DirtyThresholdRatio: 0.2,               // Dirty ratio threshold
}

cacheStore, err := store.NewCacheStore(cfg)
```

### Configuration Table

| Option                 | Description                         | Default    |
|------------------------|-------------------------------------|------------|
| `GCInterval`           | Expired key cleanup interval        | 10s        |
| `DBSave`               | Enable SQLite persistence           | true       |
| `DBFileName`           | SQLite file path                    | "cache.db" |
| `DBSaveInterval`       | Automatic DB save interval          | 10m        |
| `SaveDirtyData`        | Enable change tracking              | true       |
| `DirtyThresholdCount`  | Full sync trigger count             | 50         |
| `DirtyThresholdRatio`  | Full sync trigger ratio             | 0.2        |

## 🔧 Supported Types & Methods

| **Type**           | **Set Method**              | **Get Method**              |
|--------------------|-----------------------------|-----------------------------|
| Raw                | `SetRaw`                    | `GetRaw`                    |
| String             | `SetString`                 | `GetString`                 |
| Boolean            | `SetBool`                   | `GetBool`                   |
| Integer            | `SetInt16, SetInt32, ...`   | `GetInt16, GetInt32, ...`   |
| Unsigned Integer   | `SetUint16, SetUint32, ...` | `GetUint16, GetUint32, ...` |
| Float              | `SetFloat32, SetFloat64`    | `GetFloat32, GetFloat64`    |
| Time               | `SetTime`                   | `GetTime`                   |
| JSON               | `SetJSON`                   | `GetJSON(key, &target)`     |

## 🎯 Advanced Features

### TTL Management
```go
// Check TTL
ttl := cacheStore.TTL("key")
switch ttl {
case store.TTLNoExpiry:
    fmt.Println("No expiry")
case store.TTLExpired:
    fmt.Println("Key expired or not found")
default:
    fmt.Printf("Remaining: %v\n", ttl)
}
```

### Key Management
```go
// Get all keys
allKeys := cacheStore.Keys()

// Check existence
count := cacheStore.Exists("key1", "key2", "key3")
fmt.Printf("%d keys exist\n", count)

// Flush all data
cacheStore.Flush()
```

### Sync
```go
// Manual sync (dirty data only)
cacheStore.Sync()

// Full sync
cacheStore.FullSync()
```

### Performance Optimization: Zero-Copy
```go
// ⚠️ Warning: Don't modify returned value
dataType, value, err := cacheStore.GetNoCopy("key")
// Only read from value!
```

## 🧪 Testing

```bash
# Coverage
go test -cover ./...

# Benchmark
go test -bench=. -benchmem ./store
```

---

⭐ If you find this project helpful, please give it a star!