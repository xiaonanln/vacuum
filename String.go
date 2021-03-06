package vacuum

import (
	"fmt"

	"runtime"

	"sync/atomic"

	"time"

	"sync"

	"github.com/xiaonanln/goSyncQueue"
	"github.com/xiaonanln/goTimer"
	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
	"github.com/xiaonanln/vacuum/vlog"
)

const (
	SS_MIGRATING   = uint64(1 << iota)
	SS_FINIALIZING = uint64(1 << iota)
)

type IString interface {
	Init()
	OnReady()
	Loop(msg StringMessage)
	Fini()

	// Migrate ...
	OnMigrateOut(extraInfo map[string]interface{})
	OnMigrateIn(extraInfo map[string]interface{})
	OnMigratedAway()

	IsPersistent() bool
	GetPersistentData() map[string]interface{}
	LoadPersistentData(data map[string]interface{})

	Save()
}

type String struct {
	ID   string
	Name string

	I          IString
	initArgs   []interface{}
	inputQueue sync_queue.SyncQueue
	outputSid  string

	_flags uint64

	migratingToServerID      int
	migratingTowardsStringID string
	migrateNotify            chan int

	// timer management
	timers     map[*timer.Timer]bool
	timersLock sync.Mutex
}

func setupString(s *String, stringID string, name string, initArgs []interface{}) {
	s.ID = stringID
	s.Name = name
	s.initArgs = initArgs
	s.inputQueue = sync_queue.NewSyncQueue()
}

func (s *String) Init() {
	vlog.Debug("%s.Init ..., args=%v", s, s.Args())
}

func (s *String) Fini() {
	vlog.Debug("%s.Fini ...", s)
}

func (s *String) OnReady() {
	vlog.Debug("%s.OnReady ...", s)
}

func (s *String) OnMigrateIn(extraInfo map[string]interface{}) {

}

func (s *String) OnMigrateOut(extraInfo map[string]interface{}) {

}

func (s *String) OnMigratedAway() {
	vlog.Debug("%s.OnMigratedAway ...", s)
}

func (s *String) Loop(msg StringMessage) {
	vlog.Debug("%s.Loop %v", s, msg)
}

func (s *String) IsPersistent() bool {
	return false
}

func (s *String) GetPersistentData() map[string]interface{} {
	return nil
}

func (s *String) LoadPersistentData(data map[string]interface{}) {

}

func (s *String) String() string {
	if s != nil {
		var pp string
		if s.I.IsPersistent() {
			pp = "P!"
		} else {
			pp = ""
		}
		return fmt.Sprintf("%s<%s>%s", s.Name, s.ID, pp)
	} else {
		return "Nil<nil>"
	}
}

func (s *String) Args() []interface{} {
	return s.initArgs
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

func (s *String) AddCallback(d time.Duration, callback func()) {
	t := timer.AddCallback(d, func() {
		vlog.Info("BEFORE INPUTQUEUE PUSH")
		s.inputQueue.Push(_TaggedMsg{callback, _MSG_TAG_CALLFUNC})
		vlog.Info("AFTER INPUTQUEUE PUSH")
	})
	s.timers[t] = true
}

//func (s *String) Read() StringMessage {
//	return <-s.inputChan
//}

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

// Destroy string on local server
func (s *String) destroyString() {
	delString(s.ID) // delete the string on local server
	s.inputQueue.Close()

	// clear timers
	for t, _ := range s.timers {

	}
}

func Send(stringID string, msg interface{}) {
	if stringID == "" {
		vlog.TraceError("Send: stringID is empty")
	}

	vlog.Debug(">>> SEND %s: %v", stringID, msg)
	dispatcher_client.SendStringMessage(stringID, msg, 0)
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
