package vacuum

import (
	"fmt"

	"runtime"

	"sync/atomic"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	STRING_MESSAGE_BUFFER_SIZE = 10
)

type StringDelegate interface {
	Init(s *String, args ...interface{})
	Fini(s *String)
	Loop(s *String, msg StringMessage)
}

type StringDelegateMaker func() StringDelegate

const (
	SS_MIGRATING = uint64(1 << iota)
	SS_MIGRATING_
	SS_FINIALIZING = uint64(1 << iota)
)

type String struct {
	ID          string
	Name        string
	delegate    StringDelegate
	persistence PersistentString
	inputChan   chan StringMessage
	outputSid   string

	_flags uint64

	migratingToServerID int
	migrateNotify       chan int
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
	s.setflags(s.flags() | flag)
}

func (s *String) HasFlag(flag uint64) bool {
	return (s.flags() & flag) != 0
}

func (s *String) setflags(val uint64) {
	atomic.StoreUint64(&s._flags, val)
}

func (s *String) flags() uint64 {
	return atomic.LoadUint64(&s._flags)
}

func (s *String) IsPersistent() bool {
	return s.persistence != nil
}

func (s *String) Read() StringMessage {
	return <-s.inputChan
}

func (s *String) tryRead() (StringMessage, bool) {
	select {
	case msg, _ := <-s.inputChan:
		return msg, true
	default:
		return nil, false
	}
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

	vlog.Debug(">>> SEND %s: %v", stringID, msg)
	dispatcher_client.SendStringMessage(stringID, msg)
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
