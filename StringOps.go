package vacuum

import (
	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

// OnXxxXxxx functions are called from dispatcher client and there is only
// one dispatcher client, so there is no concurrency problem in these functions

// CreateString: create a string with specified name
func CreateString(name string, args ...interface{}) string {
	stringID := uuid.GenUUID()
	dispatcher_client.SendCreateStringReq(name, stringID, args)
	return stringID
}

// Create a String with specified name on the local vacuum server
func CreateStringLocally(name string, args ...interface{}) string {
	stringID := uuid.GenUUID()
	OnCreateString(name, stringID, args)
	dispatcher_client.SendCreateStringLocallyReq(name, stringID)
	return stringID
}

func LoadString(name string, stringID string) {
	// load string from storage
	dispatcher_client.SendLoadStringReq(name, stringID)
	//data, err := stringStorage.Read(stringID)
	//if err != nil {
	//	// load failed
	//	panic(err)
	//}
}

// OnCreateString: called when dispatcher sends create string resp
func OnCreateString(name string, stringID string, args []interface{}) {
	createString(name, stringID, args, true)
}

func OnLoadString(name string, stringID string) {
	vlog.Debugf("OnLoadString: name=%s, stringID=%s", name, stringID)
	createString(name, stringID, []interface{}{}, false)
}

func createString(name string, stringID string, args []interface{}, isNewString bool) {
	delegateMaker := getStringDelegateMaker(name)
	if delegateMaker == nil {
		vlog.Panicf("OnCreateString: routine of String %s is nil", name)
	}

	delegate := delegateMaker()
	s := newString(stringID, name, delegate)
	putString(s)
	vlog.Debugf("OnCreateString %s: %s, args=%v", name, s, args)

	go func() {
		defer onStringRoutineQuit(name, stringID)

		s.delegate.Init(s, args...)

		if !isNewString {
			data, err := stringStorage.Read(name, stringID)
			if err != nil {
				// load string failed..
				vlog.Panic(err)
			}
			if data != nil {
				s.persistence.LoadPersistentData(data.(map[string]interface{}))
			}
		}

		for {
			msg := s.Read()
			if msg != nil {
				s.delegate.Loop(s, msg)
			} else {
				break
			}
		}
		s.delegate.Fini(s)
	}()
}

// DeclareService: declare that the specified String provides specified service
func DeclareService(sid string, serviceName string) {
	dispatcher_client.SendDeclareServiceReq(sid, serviceName)
}

func OnDeclareService(stringID string, serviceName string) {
	vlog.Infof("vacuum: OnDeclareService: %s => %s", stringID, serviceName)
	declareService(stringID, serviceName)
}

func OnSendStringMessage(stringID string, msg common.StringMessage) {
	s := getString(stringID)
	vlog.Debugf("vacuum: OnSendStringMessage: %s => %v", s, msg)
	s.inputChan <- msg
}

//// Close specified string
//func Close(stringID string) {
//	s := getString(stringID)
//	if s == nil {
//		dispatcher_client.RelayCloseString(stringID)
//	} else {
//		s.Close()
//	}
//}
//
func OnCloseString(stringID string) {
	//s := getString(stringID)
	//s.Close()
}

// Called after string quit its routine
func onStringRoutineQuit(name string, stringID string) {
	vlog.Debugf("String %s.%s quited.", name, stringID)
	delString(stringID) // delete the string on local server
	undeclareServicesOfString(stringID)
	dispatcher_client.SendStringDelReq(stringID)
}

// string del notification from dispatcher
func OnDelString(stringID string) {
	undeclareServicesOfString(stringID)
}
