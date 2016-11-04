package vacuum_server

import (
	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum/config"
	"github.com/xiaonanln/vacuum/storage"
	"github.com/xiaonanln/vacuum/storage/backend/filesystem"
)

func openStorage(config *config.StorageConfig) storage.StringStorage {
	if config.Type == "filesystem" {
		return filesystem.OpenDirectory(config.Directory)
	} else {
		logrus.Panicf("openStorage: unknown storage type in conf: %s", config.Type)
		return nil
	}
}
