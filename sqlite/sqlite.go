package sqlite

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	db  *sql.DB
	mux sync.Mutex
}

func initDB(filename string) (*sql.DB, error) {
	if filename == "" {
		return nil, errors.ErrFileNameEmpty
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cache_data (
		key TEXT PRIMARY KEY,
		data_type INTEGER,
		data BLOB,
		expiry INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewSqliteStore(filename string) (*SqliteStore, error) {
	db, err := initDB(filename)
	if err != nil {
		return nil, err
	}
	return &SqliteStore{
		db: db,
	}, nil
}

func (s *SqliteStore) LoadFromDB() (map[string]entry.Entry, error) {
	if s.db == nil {
		return nil, errors.ErrDBNotInit
	}

	rows, err := s.db.Query("SELECT key, data_type, data, expiry FROM cache_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbData := make(map[string]entry.Entry)
	now := uint32(time.Now().Unix())
	for rows.Next() {
		var key string
		var dataType types.DataType
		var data []byte
		var expiry uint32

		if err := rows.Scan(&key, &dataType, &data, &expiry); err != nil {
			log.Println(err)
			continue
		}

		if expiry > 0 && expiry <= now {
			continue
		}

		dbData[key] = entry.Entry{
			Type:   dataType,
			Data:   data,
			Expiry: expiry,
		}
	}

	return dbData, nil
}

func (s *SqliteStore) SaveDirtyData(set_dirtys map[string]entry.Entry, delete_dirtys []string) error {
	if s.db == nil {
		return errors.ErrDBNotInit
	}

	if len(set_dirtys) == 0 && len(delete_dirtys) == 0 {
		return nil
	}

	if s.mux.TryLock() {
		defer s.mux.Unlock()
	} else {
		return errors.ErrAlreadySave
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	insertStmt, err := tx.Prepare(`
		INSERT INTO cache_data (key, data_type, data, expiry) 
		VALUES (?, ?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			data_type = excluded.data_type,
			data = excluded.data,
			expiry = excluded.expiry
	`)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	deleteStmt, err := tx.Prepare("DELETE FROM cache_data WHERE key = ?")
	if err != nil {
		return err
	}
	defer deleteStmt.Close()

	now := uint32(time.Now().Unix())

	for key, entry := range set_dirtys {
		if entry.IsExpiredWithTime(now) {
			continue
		}

		if _, err := insertStmt.Exec(key, entry.Type, entry.Data, entry.Expiry); err != nil {
			return err
		}
	}

	for _, key := range delete_dirtys {
		if _, err := deleteStmt.Exec(key); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SqliteStore) Save(data map[string]entry.Entry, force bool) error {
	if s.db == nil {
		return errors.ErrDBNotInit
	}
	if force {
		s.mux.Lock()
		defer s.mux.Unlock()
	} else if s.mux.TryLock() {
		defer s.mux.Unlock()
	} else {
		return errors.ErrAlreadySave
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM cache_data"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO cache_data (key, data_type, data, expiry) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := uint32(time.Now().Unix())

	for key, entry := range data {
		if entry.IsExpiredWithTime(now) {
			continue
		}

		if _, err := stmt.Exec(key, entry.Type, entry.Data, entry.Expiry); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SqliteStore) Close() error {
	if s.db == nil {
		return errors.ErrDBNotInit
	}
	return s.db.Close()
}
