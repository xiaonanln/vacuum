package concurrent

import "github.com/xiaonanln/vacuum/storage"

type ConcurrentStringStorage struct {
	subStorages []storage.StringStorage
}

func NewConcurrentStringStorage(subStorages []storage.StringStorage) storage.StringStorage {
	return &ConcurrentStringStorage{
		subStorages: subStorages,
	}
}

func hash(stringID string) uint {
	var h uint = 1
	for _, c := range []byte(stringID) {
		h = h * uint(c)
	}
	return h
}

func (ss *ConcurrentStringStorage) subStorageOf(name string, stringID string) storage.StringStorage {
	return ss.subStorages[hash(stringID)%uint(len(ss.subStorages))]
}

func (ss *ConcurrentStringStorage) Write(name string, stringID string, data interface{}) error {
	return ss.subStorageOf(name, stringID).Write(name, stringID, data)
}

func (ss *ConcurrentStringStorage) Read(name string, stringID string) (interface{}, error) {
	return ss.subStorageOf(name, stringID).Read(name, stringID)
}
