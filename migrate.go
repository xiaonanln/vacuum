package vacuum

import (
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

type Migratable interface {
	GetMigrateData(s *String) interface{}
	InitWithMigrateData(data interface{})
}

// Migrate to other server
func (s *String) Migrate(serverID int) {
	// remove string from this vacuum server
	// send the start-migration notification to dispatcher
	// migrate the data of string to vacuum server

	if s.HasFlag(SS_FINIALIZING) {
		vlog.Panicf("Do not migrate when finializing")
	}

	if s.HasFlag(SS_MIGRATING) { // already migrating...
		return
	}

	var data map[string]interface{}

	// get migrate data from string
	if s.persistence != nil {
		data = s.persistence.GetPersistentData()
	}

	s.SetFlag(SS_MIGRATING)
	popString(s.ID) // pop self from Vacuum
	dispatcher_client.SendMigrateStringReq(s.Name, s.ID, serverID, data)
}
