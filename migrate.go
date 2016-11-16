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
		vlog.Debugf("%s: already migrating", s)
		return
	}

	vlog.Debugf("%s.Migrate: start migrating ...", s)
	// mark as migrating
	s.SetFlag(SS_MIGRATING)
	// send the start-migrate req
	dispatcher_client.SendStartMigrateStringReq(s.ID)
}

func StartMigrateString(stringID string) {
	s := popString(stringID) // get the migrating string
	vlog.Debugf(">>> StartMigrateString: stringID=%s, string=%v, migrating=%v", stringID, s, s != nil && s.HasFlag(SS_MIGRATING))

	if s == nil || s.HasFlag(SS_FINIALIZING) {
		// String gone or finializing, migrate stop.
		vlog.Debugf("StartMigrateString: String %s already finialized or qutied", stringID)
		return
	}

	// whenever we started migrating,
	// there should be no more msg to this String,
	// so there should be no conflict on String

	var data map[string]interface{}

	// get migrate data from string
	if s.persistence != nil {
		data = s.persistence.GetPersistentData()
	}

	dispatcher_client.SendMigrateStringReq(s.Name, s.ID, serverID, data)
}

// String migrated to this server
func OnMigrateString(name string, stringID string, data map[string]interface{}) {
	vlog.Debugf("String %s.%s migrated to server %v: data=%v", name, stringID, serverID, data)
	createString(name, stringID, nil, false, data)
}
