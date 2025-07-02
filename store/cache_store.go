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
		if cfg.SaveDirtyData {
			if cfg.DirtyThresholdCount <= 0 {
				return nil, errors.ErrDirtyThresholdCount
			}
			if cfg.DirtyThresholdRatio <= 0 || cfg.DirtyThresholdRatio > 1 {
				return nil, errors.ErrDirtyThresholdRatio
			}
			store.dirty = newDirtyManager(cfg.DirtyThresholdCount, cfg.DirtyThresholdRatio)
		}
	} else {
		cfg.DBSaveInterval = 0
	}

	if intervalFunc := store.createTicker(cfg.GCInterval, cfg.DBSaveInterval); intervalFunc != nil {
		store.wg.Add(1)
		go func() {
			defer store.wg.Done()
			intervalFunc()
		}()
	}

	return store, nil
}
