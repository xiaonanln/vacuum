package concurrent

import (
	"sync"

	"github.com/xiaonanln/vacuum/storage"
)

type ConcurrentStringStorage struct {
	subStorages []storage.StringStorage
	locks       []sync.RWMutex
}

func NewConcurrentStringStorage(subStorages []storage.StringStorage) storage.StringStorage {
	locks := make([]sync.RWMutex, len(subStorages), len(subStorages))
	return &ConcurrentStringStorage{
		subStorages: subStorages,
		locks:       locks,
	}
}

func hash(stringID string) uint {
	var h uint = 1
	for _, c := range []byte(stringID) {
		h = h * uint(c)
	}
	return h
}

func (ss *ConcurrentStringStorage) subStorageOf(name string, stringID string) uint {
	return hash(stringID) % uint(len(ss.subStorages))
}

func (ss *ConcurrentStringStorage) Write(name string, stringID string, data interface{}) error {
	idx := ss.subStorageOf(name, stringID)
	lock := ss.locks[idx]
	lock.Lock()
	defer lock.Unlock()

	return ss.subStorages[idx].Write(name, stringID, data)
}

func (ss *ConcurrentStringStorage) Read(name string, stringID string) (interface{}, error) {
	idx := ss.subStorageOf(name, stringID)
	lock := ss.locks[idx]
	lock.RLock()
	defer lock.RUnlock()

	return ss.subStorages[idx].Read(name, stringID)
}
