package vacuum_server

import (
	"strings"

	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/storage"
	"github.com/xiaonanln/vacuum/storage/backend/concurrent"
	"github.com/xiaonanln/vacuum/storage/backend/filesystem"
	"github.com/xiaonanln/vacuum/storage/backend/mongodb"
	"github.com/xiaonanln/vacuum/vlog"
)

func openStorage(config *config.StorageConfig) storage.StringStorage {
	storageType := strings.ToLower(config.Type)
	var ss storage.StringStorage
	var err error

	var opener func() (storage.StringStorage, error)

	if storageType == "filesystem" {
		opener = func() (storage.StringStorage, error) {
			return string_storage_filesystem.OpenDirectory(config.Directory)
		}
	} else if storageType == "mongodb" {
		opener = func() (storage.StringStorage, error) {
			return string_storage_mongodb.OpenMongoDB(config.Url, config.DB)
		}
	} else {
		vlog.Panicf("openStorage: unknown storage type in conf: %s", config.Type)
	}

	if config.Concurrent <= 1 {
		ss, err = opener()
	} else {
		subStorages := make([]storage.StringStorage, config.Concurrent, config.Concurrent)
		for i := 0; i < config.Concurrent; i++ {
			subss, err := opener()
			if err != nil {
				vlog.Panic(err)
			}
			subStorages[i] = subss
		}
		ss = concurrent.NewConcurrentStringStorage(subStorages)
	}

	if err != nil {
		vlog.Panic(err)
	}
	return ss
}
