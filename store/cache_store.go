package store

import (
	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/sqlite"
)

func NewCacheStore(cfg config.Config) (*CacheStore, error) {
	store := &CacheStore{
		memorydb: make(map[string]entry.Entry),
		done:     make(chan bool),
		config:   cfg,
	}
	if cfg.DBSave {
		if cfg.DBFileName == "" {
			cfg.DBFileName = "cache.db"
		}
		db, err := sqlite.InitDB(cfg.DBFileName)
		if err != nil {
			return nil, err
		}
		defer db.Close()
		if data, err := sqlite.LoadFromDB(db); err != nil {
			return nil, err
		} else {
			store.memorydb = data
		}
	}

	if cfg.GCInterval > 0 {
		go store.gc()
	}

	return store, nil
}
