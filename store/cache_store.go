package store

import (
	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
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
			return nil, errors.ErrFileNameEmpty
		}
		db, err := sqlite.InitDB(cfg.DBFileName)
		if err != nil {
			return nil, err
		}

		data, err := sqlite.LoadFromDB(db)
		db.Close()

		if err != nil {
			return nil, err
		}
		store.memorydb = data
	}

	if cfg.GCInterval > 0 {
		go store.gc()
	}

	return store, nil
}
