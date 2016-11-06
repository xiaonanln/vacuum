package vacuum_server

import (
	"strings"

	"log"

	"fmt"

	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/storage"
	"github.com/xiaonanln/vacuum/storage/backend/filesystem"
	"github.com/xiaonanln/vacuum/storage/backend/mongodb"
)

func openStorage(config *config.StorageConfig) storage.StringStorage {
	storageType := strings.ToLower(config.Type)
	var ss storage.StringStorage
	var err error

	if storageType == "filesystem" {
		ss, err = string_storage_filesystem.OpenDirectory(config.Directory)
	} else if storageType == "mongodb" {
		ss, err = string_storage_mongodb.OpenMongoDB(config.Url, config.DB)
	} else {
		err = fmt.Errorf("openStorage: unknown storage type in conf: %s", config.Type)
	}
	if err != nil {
		log.Panic(err)
	}
	return ss
}
