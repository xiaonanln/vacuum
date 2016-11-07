package string_storage_mongodb

import (
	"fmt"

	"github.com/xiaonanln/vacuum/storage"
	"github.com/xiaonanln/vacuum/vlog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_DB_NAME = "vacuum"
)

var (
	db *mgo.Database
)

type MongoDBStringStorge struct {
	db *mgo.Database
}

func OpenMongoDB(url string, dbname string) (storage.StringStorage, error) {
	vlog.Debug("Connecting MongoDB ...")
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	if dbname == "" {
		// if db is not specified, use default
		dbname = DEFAULT_DB_NAME
	}
	db = session.DB(dbname)
	return &MongoDBStringStorge{
		db: db,
	}, nil
}

func collectionName(name string) string {
	return fmt.Sprintf("S_%s", name)
}

func (ss *MongoDBStringStorge) Write(name string, stringID string, data interface{}) error {
	col := ss.db.C(collectionName(name))
	_, err := col.UpsertId(stringID, bson.M{
		"data": data,
	})
	return err
}

func (ss *MongoDBStringStorge) Read(name string, stringID string) (interface{}, error) {
	col := ss.db.C(collectionName(name))
	q := col.FindId(stringID)
	var doc bson.M
	err := q.One(&doc)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}(doc["data"].(bson.M)), nil
}
