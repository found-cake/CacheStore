package config

import "time"

type Config struct {
	GCInterval     time.Duration
	DBSave         bool
	DBFileName     string
	DBSaveInterval time.Duration
}

func DefaultConfig() Config {
	return Config{
		GCInterval:         10 * time.Second,
		DBSave:             true,
		DBFileName:         "cache.db",
		DBSaveInterval: 10 * time.Minute,
	}
}
