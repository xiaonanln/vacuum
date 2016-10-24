package vacuum

import (
	"fmt"

	"log"

	"runtime"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	ALWAYS_SEND_STRING_MESSAGE_THROUGH_DISPATCHER = true // only use true for debug!
	STRING_MESSAGE_BUFFER_SIZE                    = 1000
)

type StringRoutine func(*String)

type String struct {
	ID        string
	Name      string
	routine   StringRoutine
	inputChan chan StringMessage
	outputSid string
}

func newString(stringID string, name string, routine StringRoutine) *String {
	return &String{
		ID:        stringID,
		Name:      name,
		routine:   routine,
		inputChan: make(chan StringMessage, STRING_MESSAGE_BUFFER_SIZE),
		outputSid: "",
	}
}

func (s *String) String() string {
	if s != nil {
		return fmt.Sprintf("String<%s.%s>", s.Name, s.ID)
	} else {
		return "String<nil>"
	}
}

func (s *String) Read() StringMessage {
	msg, _ := <-s.inputChan
	//log.Debugf("%s READ %T(%v)", s, msg, msg)
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
		log.Panicf("Send: stringID is empty")
	}

	if !ALWAYS_SEND_STRING_MESSAGE_THROUGH_DISPATCHER {
		targetString := getString(stringID)

		if targetString == nil { // string is not local, send msg to dispatcher
			dispatcher_client.SendStringMessage(stringID, msg)
		} else { // found the target string on this vacuum server
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
