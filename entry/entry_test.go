package entry

import (
	"testing"
	"time"

	"github.com/found-cake/CacheStore/utils/types"
)

func TestNewEntry_ExpirySet(t *testing.T) {
	data := []byte("test")
	exp := 2 * time.Second
	entry := NewEntry(types.RAW, data, exp)

	if entry.Type != types.RAW {
		t.Errorf("expected Type %v, got %v", types.RAW, entry.Type)
	}
	if string(entry.Data) != "test" {
		t.Errorf("expected Data 'test', got '%s'", string(entry.Data))
	}
	if entry.Expiry == 0 {
		t.Errorf("expiry should be set for positive duration")
	}
}

func TestNewEntry_NoExpiry(t *testing.T) {
	data := []byte("test")
	entry := NewEntry(types.RAW, data, 0)

	if entry.Expiry != 0 {
		t.Errorf("expected Expiry 0 for non-expiring entry, got %d", entry.Expiry)
	}
}

func TestIsExpired(t *testing.T) {
	entry := NewEntry(types.RAW, []byte("test"), 1*time.Second)
	if entry.IsExpired() {
		t.Error("entry should not be expired immediately after creation")
	}
	time.Sleep(1100 * time.Millisecond)
	if !entry.IsExpired() {
		t.Error("entry should be expired after duration")
	}
}

func TestIsExpiredWithTime(t *testing.T) {
	now := time.Now().Unix()
	entry := Entry{
		Type:   types.STRING,
		Data:   []byte("test"),
		Expiry: now + 1,
	}
	if entry.IsExpiredWithTime(now) {
		t.Error("should not be expired at current time")
	}
	if !entry.IsExpiredWithTime(now + 2) {
		t.Error("should be expired after expiry")
	}
}
