package vacuum

import (
	"log"

	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

func CreateString(name string) {
	dispatcher_client.SendCreateStringReq(name)
}

// OnCreateString: called when dispatcher sends create string resp
func OnCreateString(name string) {
	routine := getStringRoutine(name)
	s := newString(routine)
	putString(s)
	go s.routine(s)
}

// DeclareService: declare that the specified String provides specified service
func DeclareService(sid string, serviceName string) {
	dispatcher_client.SendDeclareServiceReq(sid, serviceName)
}

func OnDeclareService(stringID string, serviceName string) {
	log.Printf("vacuum: OnDeclareService: %s => %s", stringID, serviceName)
	declareService(stringID, serviceName)
}
