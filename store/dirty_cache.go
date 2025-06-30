package store

import "sync"

type DirtyAction bool

const (
	DirtySet    DirtyAction = true
	DirtyDelete             = false
)

type dirtyManager struct {
	mux          sync.RWMutex
	dirtyData    map[string]DirtyAction
	needFullSync bool
}

func NewDirtyManager() *dirtyManager {
	return &dirtyManager{
		dirtyData: make(map[string]DirtyAction),
	}
}

func (d *dirtyManager) UnsafeClear() {
	d.dirtyData = make(map[string]DirtyAction)
}

func (d *dirtyManager) Clear() {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.dirtyData = make(map[string]DirtyAction)
}

func (d *dirtyManager) UnsafeSet(key string) {
	if !d.needFullSync {
		d.dirtyData[key] = DirtySet
	}
}

func (d *dirtyManager) Set(key string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.UnsafeSet(key)
}

func (d *dirtyManager) UnsafeDelete(key string) {
	if !d.needFullSync {
		d.dirtyData[key] = DirtyDelete
	}
}

func (d *dirtyManager) Delete(key string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.UnsafeDelete(key)
}

func (d *dirtyManager) Size() int {
	return len(d.dirtyData)
}

func (d *dirtyManager) Keys() ([]string, []string) {
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
