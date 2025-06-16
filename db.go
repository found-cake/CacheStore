package cachestore

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func initDB(filename string) (*sql.DB, error) {
	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cache_data (
		key TEXT PRIMARY KEY,
		data BLOB,
		expiry INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadFromDB(db *sql.DB) (map[string]entry, error) {
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	rows, err := db.Query("SELECT key, data, expiry FROM cache_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbData := make(map[string]entry)
	now := uint32(time.Now().Unix())
	for rows.Next() {
		var key string
		var data []byte
		var expiry uint32

		if err := rows.Scan(&key, &data, &expiry); err != nil {
			log.Println(err)
			continue
		}

		if expiry > 0 && expiry <= now {
			continue
		}

		dbData[key] = entry{
			data:   data,
			expiry: expiry,
		}
	}

	return dbData, nil
}

func saveDB(db *sql.DB, data map[string]entry) error {
	if db == nil {
		return errors.New("database not initialized")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM cache_data"); err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO cache_data (key, data, expiry) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for key, entry := range data {
		if entry.expiry > 0 && entry.expiry <= uint32(time.Now().Unix()) {
			continue
		}

		if _, err := stmt.Exec(key, entry.data, entry.expiry); err != nil {
			return err
		}
	}

	return tx.Commit()
}
