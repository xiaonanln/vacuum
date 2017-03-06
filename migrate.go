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
	s.startMigrate(serverID, "")
}

func (s *String) MigrateTowards(otherID string) {
	vlog.Debug("MigrateTowards %s", otherID)
	s.startMigrate(0, otherID)
}

func (s *String) startMigrate(serverID int, otherID string) {
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
	if serverID > 0 {
		WaitServerReady(serverID)
	}

	vlog.Debug("%s.Migrate: start migrating ...", s)
	// mark as migrating
	s.migratingToServerID = serverID
	s.migratingTowardsStringID = otherID
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
func OnMigrateString(name string, stringID string, initArgs []interface{}, data map[string]interface{}, extraMigrateInfo map[string]interface{}) {
	vlog.Debug("String %s.%s migrated to server %v: args=%v, data=%v, extraMigrateInfo=%v", name, stringID, serverID, initArgs, data, extraMigrateInfo)
	createString(name, stringID, initArgs, CREATE_MIGRATE, data, extraMigrateInfo)
}
