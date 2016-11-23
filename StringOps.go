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
	createString(name, stringID, args, false, nil)
}

func OnLoadString(name string, stringID string) {
	vlog.Debug("OnLoadString: name=%s, stringID=%s", name, stringID)
	createString(name, stringID, []interface{}{}, true, nil)
}

func createString(name string, stringID string, args []interface{}, loadFromStorage bool, migrateData map[string]interface{}) {
	delegateMaker := getStringDelegateMaker(name)
	if delegateMaker == nil {
		vlog.Panicf("OnCreateString: routine of String %s is nil", name)
	}

	delegate := delegateMaker()
	s := newString(stringID, name, args, delegate)
	putString(s)
	vlog.Debug("OnCreateString %s: %s, args=%v", name, s, args)

	go func() {
		defer recoverFromStringRoutineError(s)

		s.delegate.Init(s, args...)

		if loadFromStorage { // loading string from storage ...
			data, err := stringStorage.Read(name, stringID)
			if err != nil {
				// load string failed..
				vlog.Panic(err)
			}
			if data != nil {
				s.persistence.LoadPersistentData(data.(map[string]interface{}))
			}
		} else if migrateData != nil { // migrated from from other server ...
			s.persistence.LoadPersistentData(migrateData)
		} else { // creating new string
			if s.IsPersistent() { // save persistent string right after it's created && inited
				s.Save()
			}
		}

		for {
			if s.HasFlag(SS_MIGRATING) {
				goto migrating_wait_notify
			}

			msg := s.inputQueue.Pop()
			if msg != nil {
				s.delegate.Loop(s, msg)
			} else {
				s.SetFlag(SS_FINIALIZING)
				goto finialize_string
			}

		}

	finialize_string:
		if s.HasFlag(SS_MIGRATING) {
			vlog.Debug("%s: string migrated ignored because it's finializing", s)
		}

		s.delegate.Fini(s)

		if s.IsPersistent() {
			s.Save()
		}

		vlog.Debug("--- %s quited", s)
		onStringRoutineQuit(s)
		return

	migrating_wait_notify:
		vlog.Debug("%s: Waiting for StartMigrateStringResp from dispatcher ...", s)
		<-s.migrateNotify
		vlog.Debug("%s: StartMigrateStringResp OK")
		// process all pending messages
	migrating_read_loop:
		for {
			msg, ok := s.inputQueue.TryPop()
			if ok {
				if msg != nil {
					s.delegate.Loop(s, msg)
				} else {
					s.SetFlag(SS_FINIALIZING)
					break migrating_read_loop
				}
			} else {
				// no more messages, now we can quit
				break migrating_read_loop
			}
		}
		// all messages are processed, now we can start migrate or quit
		vlog.Debug("%s: all messages handled, quiting for migrating ... finializing=%v", s, s.HasFlag(SS_FINIALIZING))
		if s.HasFlag(SS_FINIALIZING) {
			// if the string is finializing, migrating is shutdown
			goto finialize_string
		}
		// real migrate now!
		var data map[string]interface{}

		// get migrate data from string
		if s.persistence != nil {
			data = s.persistence.GetPersistentData()
		}

		dispatcher_client.SendMigrateStringReq(s.Name, s.ID, s.migratingToServerID, s.initArgs, data)
		return
	}()
}

// DeclareService: declare that the specified String provides specified service
func DeclareService(sid string, serviceName string) {
	dispatcher_client.SendDeclareServiceReq(sid, serviceName)
}

func OnDeclareService(stringID string, serviceName string) {
	declareService(stringID, serviceName)
}

func OnSendStringMessage(stringID string, msg common.StringMessage) {
	s := getString(stringID)
	vlog.Debug("vacuum: OnSendStringMessage: %s: %s => %v", stringID, s, msg)
	s.inputQueue.Push(msg)
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
func onStringRoutineQuit(s *String) {
	s.inputQueue.Close()
	stringID := s.ID
	delString(stringID) // delete the string on local server
	dispatcher_client.SendStringDelReq(stringID)
}

func recoverFromStringRoutineError(s *String) {
	if err := recover(); err != nil {
		vlog.TraceError("!!! %s paniced: %v", s, err)
	}
}

// string del notification from dispatcher
func OnDelString(stringID string) {
	undeclareServicesOfString(stringID)
}
