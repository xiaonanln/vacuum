package vacuum

import (
	"fmt"

	"runtime"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	ALWAYS_SEND_STRING_MESSAGE_THROUGH_DISPATCHER = true // only use true for debug!
	STRING_MESSAGE_BUFFER_SIZE                    = 0
)

type StringDelegate interface {
	Init(s *String, args ...interface{})
	Fini(s *String)
	Loop(s *String, msg StringMessage)
}

type StringDelegateMaker func() StringDelegate

const (
	SS_MIGRATING   = uint64(1 << iota)
	SS_FINIALIZING = uint64(1 << iota)
)

type String struct {
	ID          string
	Name        string
	delegate    StringDelegate
	persistence PersistentString
	inputChan   chan StringMessage
	outputSid   string
	flags       uint64
}

func newString(stringID string, name string, delegate StringDelegate) *String {
	s := &String{
		ID:        stringID,
		Name:      name,
		delegate:  delegate,
		inputChan: make(chan StringMessage, STRING_MESSAGE_BUFFER_SIZE),
		outputSid: "",
	}

	s.persistence, _ = delegate.(PersistentString)
	return s
}

func (s *String) String() string {
	if s != nil {
		var pp string
		if s.persistence != nil {
			pp = "P!"
		} else {
			pp = ""
		}
		return fmt.Sprintf("%s<%s>%s", s.Name, s.ID, pp)
	} else {
		return "Nil<nil>"
	}
}

func (s *String) SetFlag(flag uint64) {
	s.flags |= flag
}

func (s *String) HasFlag(flag uint64) bool {
	return (s.flags & flag) != 0
}

func (s *String) IsPersistent() bool {
	return s.persistence != nil
}

func (s *String) Read() StringMessage {
	msg, _ := <-s.inputChan
	//vlog.Debugf("%s READ %T(%v)", s, msg, msg)
	return msg
}

func (s *String) Output(msg StringMessage) {
	s.Send(s.outputSid, msg)
}

func (s *String) Connect(stringID string) {
	s.outputSid = stringID
}

func (s *String) Send(stringID string, msg StringMessage) {
	//dispatcher_client.SendStringMessage(sid, msg) // debug code
	Send(stringID, msg)
}

func (s *String) DeclareService(name string) {
	DeclareService(s.ID, name)
}

func (s *String) SendToService(serviceName string, msg StringMessage) {
	SendToService(serviceName, msg)
}

//// Close the String input
//func (s *String) Close() {
//	close(s.inputChan)
//}

func (s *String) Yield() {
	runtime.Gosched()
}

func Send(stringID string, msg interface{}) {
	if stringID == "" {
		vlog.Panicf("Send: stringID is empty")
	}

	//dispatcher_client.SendStringMessage(stringID, msg)

	if !ALWAYS_SEND_STRING_MESSAGE_THROUGH_DISPATCHER {
		targetString := getString(stringID)

		if targetString == nil { // string is not local, send msg to dispatcher
			dispatcher_client.SendStringMessage(stringID, msg)
		} else { // found the target string on this vacuum server
			vlog.Infof("Send through channel %s <- %v", targetString, msg)
			targetString.inputChan <- msg
		}
	} else {
		// FOR DEBUG ONLY
		dispatcher_client.SendStringMessage(stringID, msg)
	}
}

func SendToService(serviceName string, msg StringMessage) {
	stringID := chooseServiceString(serviceName)
	Send(stringID, msg)
}

func BroadcastToService(serviceName string, msg StringMessage) {
	stringIDs := getAllServiceStringIDs(serviceName)
	for _, stringID := range stringIDs {
		Send(stringID, msg)
	}
}
