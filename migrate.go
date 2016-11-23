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
		vlog.Debug("%s: already migrating", s)
		return
	}

	// wait the target server ready before migrating ...
	WaitServerReady(serverID)

	vlog.Debug("%s.Migrate: start migrating ...", s)
	// mark as migrating
	s.migratingToServerID = serverID
	s.migrateNotify = make(chan int, 1)
	s.SetFlag(SS_MIGRATING)
	// send the start-migrate req
	dispatcher_client.SendStartMigrateStringReq(s.ID)
}

func MigrateString(stringID string) {
	s := popString(stringID) // get the migrating string

	vlog.Debug(">>> MigrateString: stringID=%s, string=%v, migrating=%v", stringID, s, s != nil && s.HasFlag(SS_MIGRATING))
	if s == nil || s.HasFlag(SS_FINIALIZING) || !s.HasFlag(SS_MIGRATING) {
		// String gone or finializing, migrate stop.
		vlog.Debug("MigrateString: String %s already finialized or quited", stringID)
		return
	}

	vlog.Debug("MigrateString: transfering to string ...")
	s.migrateNotify <- 1

	vlog.Debug("MigrateString: waiting for string to complete ...")
}

// String migrated to this server
func OnMigrateString(name string, stringID string, data map[string]interface{}) {
	vlog.Debug("String %s.%s migrated to server %v: data=%v", name, stringID, serverID, data)
	createString(name, stringID, nil, false, data)
}
