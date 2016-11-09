package vacuum

import "github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"

type Migratable interface {
	GetMigrateData(s *String) interface{}
	InitWithMigrateData(data interface{})
}

// Migrate to other server
func (s *String) Migrate(serverID int) {
	// remove string from this vacuum server
	// send the start-migration notification to dispatcher
	// migrate the data of string to vacuum server

	var data map[string]interface{}
	if s.persistence != nil {
		data = s.persistence.GetPersistentData()
	}
	// get migrate data from string
	dispatcher_client.SendMigrateStringReq(s.Name, s.ID, serverID, data)
}
