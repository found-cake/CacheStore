package cachestore

import "time"

type Config struct {
	GCInterval time.Duration
	DBSave     bool
	DBFileName string
}

func DefaultConfig() Config {
	return Config{
		GCInterval: 10 * time.Second,
		DBSave:     true,
		DBFileName: "cache.db",
	}
}
