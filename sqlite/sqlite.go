package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/store/types"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filename string) (*sql.DB, error) {
	if filename == "" {
		return nil, errors.ErrFileNameEmpty
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

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

func LoadFromDB(db *sql.DB) (map[string]entry.Entry, error) {
	if db == nil {
		return nil, errors.ErrDBNotInit
	}

	rows, err := db.Query("SELECT key, data_type, data, expiry FROM cache_data")
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
			Data:   data,
			Expiry: expiry,
		}
	}

	return dbData, nil
}

func SaveDB(db *sql.DB, data map[string]entry.Entry) error {
	if db == nil {
		return errors.ErrDBNotInit
	}

	tx, err := db.Begin()
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
