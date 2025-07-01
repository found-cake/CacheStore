package store

import "sync"

type DirtyAction bool

const (
	DirtySet    DirtyAction = true
	DirtyDelete             = false
)

type dirtyManager struct {
	mux            sync.RWMutex
	dirtyData      map[string]DirtyAction
	needFullSync   bool
	ThresholdCount int
	ThresholdRatio float64
}

func newDirtyManager(count int, ratio float64) *dirtyManager {
	return &dirtyManager{
		dirtyData:      make(map[string]DirtyAction),
		ThresholdCount: count,
		ThresholdRatio: ratio,
	}
}

func (d *dirtyManager) unsafeClear() {
	d.dirtyData = make(map[string]DirtyAction)
}

func (d *dirtyManager) clear() {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.dirtyData = make(map[string]DirtyAction)
}

func (d *dirtyManager) unsafeSet(key string) {
	if !d.needFullSync {
		d.dirtyData[key] = DirtySet
	}
}

func (d *dirtyManager) set(key string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.unsafeSet(key)
}

func (d *dirtyManager) unsafeDelete(key string) {
	if !d.needFullSync {
		d.dirtyData[key] = DirtyDelete
	}
}

func (d *dirtyManager) delete(key string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.unsafeDelete(key)
}

func (d *dirtyManager) size() int {
	return len(d.dirtyData)
}

func (d *dirtyManager) keys() ([]string, []string) {
	set_keys := make([]string, 0, len(d.dirtyData))
	delete_keys := make([]string, 0, len(d.dirtyData))
	for key, action := range d.dirtyData {
		if action {
			set_keys = append(set_keys, key)
		} else {
			delete_keys = append(delete_keys, key)
		}
	}

	return set_keys, delete_keys
}

func (d *dirtyManager) NeedFullSync() {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.needFullSync = true
	d.dirtyData = make(map[string]DirtyAction)
}
