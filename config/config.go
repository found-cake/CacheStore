package config

import "time"

type Config struct {
	GCInterval          time.Duration
	DBSave              bool
	DBFileName          string
	DBSaveInterval      time.Duration
	SaveDirtyData       bool
	DirtyThresholdCount int
	DirtyThresholdRatio float64
}

func DefaultConfig() Config {
	return Config{
		GCInterval:          10 * time.Second,
		DBSave:              true,
		DBFileName:          "cache.db",
		DBSaveInterval:      10 * time.Minute,
		SaveDirtyData:       true,
		DirtyThresholdCount: 50,
		DirtyThresholdRatio: 0.2,
	}
}
