package vacuum

import "github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"

type Migratable interface {
	GetMigrateData(s *String) interface{}
	InitWithMigrateData(data interface{})
}

// Migrate to other server
func (s *String) Migrate(serverID int) {
	Migrate(s.ID, serverID)
	// remove string from this vacuum server
	// send the start-migration notification to dispatcher
	// migrate the data of string to vacuum server

	//
}

func Migrate(stringID string, serverID int) {
	_ = popString(stringID)
	// get migrate data from string
	dispatcher_client.SendMigrateStringReq(stringID, serverID)
}
