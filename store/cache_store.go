package store

import (
	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/sqlite"
)

func NewCacheStore(cfg config.Config) (*CacheStore, error) {
	store := &CacheStore{
		memorydb: make(map[string]entry.Entry),
		done:     make(chan struct{}),
	}
	if cfg.DBSave {
		sqlitedb, err := sqlite.NewSqliteStore(cfg.DBFileName)
		if err != nil {
			return nil, err
		}
		data, err := sqlitedb.LoadFromDB()
		if err != nil {
			return nil, err
		}
		store.memorydb = data
		store.sqlitedb = sqlitedb
	} else {
		cfg.DBSaveInterval = 0
	}

	if intervalFunc := store.createTicker(cfg.GCInterval, cfg.DBSaveInterval); intervalFunc != nil {
		go intervalFunc()
	}

	return store, nil
}
