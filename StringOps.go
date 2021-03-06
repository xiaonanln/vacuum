package vacuum

import (
	"reflect"

	"github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/uuid"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	_MSG_TAG_NORMAL   = 0
	_MSG_TAG_CALLFUNC = 1
)

type _TaggedMsg struct {
	msg common.StringMessage
	tag uint
}

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

func LoadString(name string, stringID string, args ...interface{}) {
	dispatcher_client.SendLoadStringReq(name, stringID, args)
}

// OnCreateString: called when dispatcher sends create string resp
func OnCreateString(name string, stringID string, args []interface{}) {
	createString(name, stringID, args, _CREATE_NEW, nil, nil)
}

func OnLoadString(name string, stringID string, args []interface{}) {
	vlog.Debug("OnLoadString: name=%s, stringID=%s, args=%v", name, stringID, args)
	createString(name, stringID, args, _CREATE_LOAD, nil, nil)
}

const ( // different ways of create string
	_CREATE_NEW     = iota + 1
	_CREATE_LOAD    = iota + 1
	_CREATE_MIGRATE = iota + 1
)

func createString(name string, stringID string, args []interface{}, createWay int, migrateData map[string]interface{}, extraMigrateInfo map[string]interface{}) {
	stringType := getRegisteredStringType(name)
	//if stringType {
	//	vlog.Panicf("OnCreateString: String type %s is unknown", name)
	//}

	derived := reflect.New(stringType) // create the new string
	s := reflect.Indirect(derived).FieldByName("String").Addr().Interface().(*String)
	setupString(s, stringID, name, args)
	s.I = derived.Interface().(IString)
	putString(s)
	vlog.Debug("OnCreateString %s: %s, args=%v", name, s, args)

	go stringRoutine(s, createWay, migrateData, extraMigrateInfo)
}
func stringRoutine(s *String, createWay int, migrateData map[string]interface{}, extraMigrateInfo map[string]interface{}) {
	vlog.Debug("stringRoutine %v %v %v %v", migrateData, extraMigrateInfo, migrateData == nil, extraMigrateInfo == nil)
	var is IString = s.I

	defer recoverFromStringRoutineError(s)

	is.Init()

	if createWay == _CREATE_LOAD { // loading string from storage ...
		data, err := stringStorage.Read(s.Name, s.ID)
		if err != nil {
			// load string failed..
			vlog.Panic(err)
		}
		if data != nil {
			is.LoadPersistentData(data.(map[string]interface{}))
		}
	} else if createWay == _CREATE_MIGRATE { // migrated from from other server ...
		is.LoadPersistentData(migrateData)
		is.OnMigrateIn(extraMigrateInfo)
	} else { // creating new string
		if is.IsPersistent() { // save persistent string right after it's created && inited
			is.Save()
		}
	}

	is.OnReady()

	for {
		if s.HasFlag(SS_MIGRATING) {
			goto migrating_wait_notify
		}

		msg := s.inputQueue.Pop().(_TaggedMsg)
		if msg.tag == _MSG_TAG_NORMAL {
			if msg.msg != nil {
				is.Loop(msg.msg)
			} else {
				s.SetFlag(SS_FINIALIZING)
				break
			}
		} else if msg.tag == _MSG_TAG_CALLFUNC {
			msg.msg.(func())()
		}
	}

finialize_string:
	vlog.Debug("%s#finialize_string ...", s)
	if s.HasFlag(SS_MIGRATING) {
		vlog.Debug("%s: string migrated ignored because it's finializing", s)
	}

	is.Fini()

	if is.IsPersistent() {
		is.Save()
	}

	vlog.Debug("--- %s quited", s)
	onStringRoutineQuit(s)
	return

migrating_wait_notify:
	vlog.Debug("%s#migrating_wait_notify: Waiting for StartMigrateStringResp from dispatcher ...", s)
	<-s.migrateNotify
	vlog.Debug("%s: StartMigrateStringResp OK", s)
	// process all pending messages

	// TODO: make strings panic safe
	for {
		_msg, ok := s.inputQueue.TryPop()
		if ok {
			msg := _msg.(_TaggedMsg)
			if msg.tag == _MSG_TAG_NORMAL {
				if msg.msg != nil {
					is.Loop(msg.msg)
				} else {
					s.SetFlag(SS_FINIALIZING)
					break
				}
			} else if msg.tag == _MSG_TAG_CALLFUNC {
				msg.msg.(func())()
			}

		} else {
			// no more messages, now we can quit
			break
		}
	}
	// all messages are processed, now we can start migrate or quit
	vlog.Debug("%s: all messages handled, quiting for migrating ... finializing=%v", s, s.HasFlag(SS_FINIALIZING))
	if s.HasFlag(SS_FINIALIZING) {
		// if the string is finializing, migrating is shutdown
		goto finialize_string
	}
	// real migrate now!
	s.destroyString() // delete the string on local server

	// get migrate data from string
	extraMigrateInfo = map[string]interface{}{}
	is.OnMigrateOut(extraMigrateInfo)
	data := is.GetPersistentData()

	dispatcher_client.SendMigrateStringReq(s.Name, s.ID, s.migratingToServerID, s.migratingTowardsStringID, s.initArgs, data, extraMigrateInfo)
	is.OnMigratedAway()
	return
}

// TODO: consider using dispatcher for service maintaince
// DeclareService: declare that the specified String provides specified service
func DeclareService(sid string, serviceName string) {
	dispatcher_client.SendDeclareServiceReq(sid, serviceName)
}

func OnDeclareService(stringID string, serviceName string) {
	declareService(stringID, serviceName)
}

func OnSendStringMessage(stringID string, msg common.StringMessage, tag uint) {
	s := getString(stringID)
	vlog.Debug("vacuum: OnSendStringMessage: %s: %s => %v, TAG %v", stringID, s, msg, tag)
	if s == nil {
		vlog.Warn("String %s not found while receiving message: %v", stringID, msg)
		return
	}

	s.inputQueue.Push(_TaggedMsg{msg, tag})
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
	stringID := s.ID
	s.destroyString() // delete the string on local server
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
