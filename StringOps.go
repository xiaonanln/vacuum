package vacuum

import (
	log "github.com/Sirupsen/logrus"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

// OnXxxXxxx functions are called from dispatcher client and there is only
// one dispatcher client, so there is no concurrency problem in these functions

// CreateString: create a string with specified name
func CreateString(name string) string {
	stringID := uuid.GenUUID()
	dispatcher_client.SendCreateStringReq(name, stringID)
	return stringID
}

// Create a String with specified name on the local vacuum server
func CreateStringLocally(name string) string {
	stringID := uuid.GenUUID()
	OnCreateString(name, stringID)
	dispatcher_client.SendCreateStringLocallyReq(name, stringID)
	return stringID
}

// OnCreateString: called when dispatcher sends create string resp
func OnCreateString(name string, stringID string) {
	routine := getStringRoutine(name)
	s := newString(stringID, name, routine)
	putString(s)
	log.Debugf("OnCreateString %s: %s", name, s)

	go s.routine(s)
}

// DeclareService: declare that the specified String provides specified service
func DeclareService(sid string, serviceName string) {
	dispatcher_client.SendDeclareServiceReq(sid, serviceName)
}

func OnDeclareService(stringID string, serviceName string) {
	log.Infof("vacuum: OnDeclareService: %s => %s", stringID, serviceName)
	declareService(stringID, serviceName)
}

func OnSendStringMessage(stringID string, msg common.StringMessage) {
	s := getString(stringID)
	log.WithField("stringID", stringID).Debugf("vacuum: OnSendStringMessage: %s => %v", s, msg)
	s.inputChan <- msg
}

// Close specified string
func Close(stringID string) {
	s := getString(stringID)
	if s == nil {
		dispatcher_client.RelayCloseString(stringID)
	} else {
		s.Close()
	}
}

func OnCloseString(stringID string) {
	s := getString(stringID)
	s.Close()
}
