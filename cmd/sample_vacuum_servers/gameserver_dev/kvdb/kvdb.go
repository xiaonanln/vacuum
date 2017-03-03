package kvdb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xiaonanln/vacuum/vlog"
)

var (
	db *leveldb.DB
)

func init() {
	var err error
	db, err = leveldb.OpenFile("game.db", nil)
	if err != nil {
		vlog.Panic(err)
		return
	}
	vlog.Info("OPEN LEVELDB OK")
}

func Get(key string, defaultVal string) string {
	v, err := db.Get([]byte(key), nil)
	if err != nil {
		//vlog.Debug("Get error: %T %v %s %v", err, err, err.Error(), err == leveldb.ErrNotFound)
		if err == leveldb.ErrNotFound {
			return defaultVal
		}
		vlog.Panic(err)
	}
	return string(v)
}

func Set(key string, val string) {
	db.Put([]byte(key), []byte(val), nil)
}
