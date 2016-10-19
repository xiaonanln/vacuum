package vacuum

import (
	"log"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

// OnXxxXxxx functions are called from dispatcher client and there is only
// one dispatcher client, so there is no concurrency problem in these functions

func CreateString(name string) string {
	stringID := uuid.GenUUID()
	dispatcher_client.SendCreateStringReq(name, stringID)
	return stringID
}

// OnCreateString: called when dispatcher sends create string resp
func OnCreateString(name string, stringID string) {
	routine := getStringRoutine(name)
	s := newString(stringID, name, routine)
	putString(s)
	log.Printf("String created: %s", s)
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

func OnSendStringMessage(stringID string, msg common.StringMessage) {
	log.Printf("vacuum: OnSendStringMessage: %s => %v", stringID, msg)
	s := getString(stringID)
	s.inputChan <- msg
}
