package sqlite

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/found-cake/CacheStore/entry"
	"github.com/found-cake/CacheStore/errors"
	"github.com/found-cake/CacheStore/utils/types"
)

var defaultData = map[string]entry.Entry{
	"foo": entry.NewEntry(types.RAW, []byte("bar"), 5*time.Minute),
}

func tempDBFile(t *testing.T) string {
	dir := t.TempDir()
	return filepath.Join(dir, "testdb.sqlite")
}

func TestNewSqliteStore_InvalidFileName(t *testing.T) {
	_, err := NewSqliteStore("")
	if err == nil {
		t.Error("expected error for empty filename")
	}
}

func TestNewSqliteStore_NoDBInit(t *testing.T) {
	s := &SqliteStore{}
	_, err := s.LoadFromDB()
	if err == nil {
		t.Error("expected error for not initialized sql")
	}
	err = s.Save(defaultData, false)
	if err == nil {
		t.Error("expected error for not initialized sql")
	}
	err = s.SaveDirtyData(defaultData, nil)
	if err == nil {
		t.Error("expected error for not initialized sql")
	}
	err = s.Close()
	if err == nil {
		t.Error("expected error for not initialized sql")
	}
}

func TestSqliteStore_Save_ForceLock(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	store.mux.Lock()
	unlocked := make(chan struct{})
	go func() {
		defer close(unlocked)
		err := store.Save(defaultData, true)
		if err != nil {
			t.Errorf("force lock Save failed: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	store.mux.Unlock()
	<-unlocked

	loaded, err := store.LoadFromDB()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got, ok := loaded["foo"]; !ok || string(got.Data) != "bar" {
		t.Errorf("expected foo to be saved, got %v", got)
	}
}

func TestSqliteStore_Save_TryLock(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	store.mux.Lock()
	defer store.mux.Unlock()

	err = store.Save(defaultData, false)
	if err != errors.ErrAlreadySave {
		t.Errorf("expected ErrAlreadySave, got %v", err)
	}
}

func TestSqliteStore_SaveDirtyData(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	err = store.Save(defaultData, true)
	if err != nil {
		t.Errorf("force lock Save failed: %v", err)
	}

	loaded, err := store.LoadFromDB()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got, ok := loaded["foo"]; !ok || string(got.Data) != "bar" {
		t.Errorf("expected foo to be saved, got %v", got)
	}

	dirtyData := map[string]entry.Entry{
		"found": entry.NewEntry(types.RAW, []byte("cake"), 5*time.Minute),
	}

	store.SaveDirtyData(dirtyData, []string{"foo"})

	loaded, err = store.LoadFromDB()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got, ok := loaded["foo"]; ok {
		t.Errorf("expected foo to be delete, got %v", got)
	}
	if got, ok := loaded["found"]; !ok || string(got.Data) != "cake" {
		t.Errorf("expected found to be saved, got %v", got)
	}
}

func TestSqliteStore_SaveDirtyData_TryLock(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	store.mux.Lock()
	defer store.mux.Unlock()

	err = store.SaveDirtyData(defaultData, nil)
	if err != errors.ErrAlreadySave {
		t.Errorf("expected ErrAlreadySave, got %v", err)
	}
}

func TestSqliteStore_Load(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	err = store.Save(defaultData, true)
	if err != nil {
		t.Errorf("force lock Save failed: %v", err)
	}
	loaded, err := store.LoadFromDB()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got, ok := loaded["foo"]; !ok || string(got.Data) != "bar" {
		t.Errorf("expected foo to be saved, got %v", got)
	}
}

func TestSqliteStore_Load_FilterExpired(t *testing.T) {
	dbfile := tempDBFile(t)
	store, err := NewSqliteStore(dbfile)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	data := map[string]entry.Entry{
		"foo": entry.NewEntry(types.RAW, []byte("bar"), 100*time.Millisecond),
	}
	err = store.Save(data, true)
	if err != nil {
		t.Errorf("force lock Save failed: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
	loaded, err := store.LoadFromDB()
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if got, ok := loaded["foo"]; ok {
		t.Errorf("expected foo to be delete, got %v", got)
	}
}
