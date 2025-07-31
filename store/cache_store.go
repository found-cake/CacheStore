package store

import (
	"time"

	"github.com/found-cake/CacheStore/config"
	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/sqlite"
)

func NewCacheStore(cfg config.Config) (*CacheStore, error) {
	store := &CacheStore{
		memorydbTemporary: make(map[string]entry.Entry),
		done:              make(chan struct{}),
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
		store.memorydbTemporary = data
		store.sqlitedb = sqlitedb
		if cfg.SaveDirtyData {
			if cfg.DirtyThresholdCount <= 0 {
				return nil, errors.ErrDirtyThresholdCount
			}
			if cfg.DirtyThresholdRatio <= 0 || cfg.DirtyThresholdRatio > 1 {
				return nil, errors.ErrDirtyThresholdRatio
			}
			store.dirty = newDirtyManager(cfg.DirtyThresholdCount, cfg.DirtyThresholdRatio)
			if cfg.DBSaveInterval > 0 {
				store.wg.Add(1)
				go func() {
					defer store.wg.Done()
					ticker := time.NewTicker(cfg.DBSaveInterval)
					defer ticker.Stop()
					for {
						select {
						case <-ticker.C:
							store.Sync()
						case <-store.done:
							return
						}
					}
				}()
			}
		}
	}

	if cfg.GCInterval > 0 {
		store.wg.Add(1)
		go func() {
			defer store.wg.Done()
			ticker := time.NewTicker(cfg.GCInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					store.cleanExpired()
				case <-store.done:
					return
				}
			}
		}()
	}

	return store, nil
}
