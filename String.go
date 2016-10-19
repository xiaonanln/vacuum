package vacuum

import (
	"fmt"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
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
		inputChan: make(chan StringMessage),
		outputSid: "",
	}
}

func (s *String) String() string {
	return fmt.Sprintf("String<%s.%s>", s.Name, s.ID)
}

func (s *String) Read() StringMessage {
	msg, _ := <-s.inputChan
	//log.Printf("%s READ %T(%v)", s, msg, msg)
	return msg
}

func (s *String) Output(msg StringMessage) {
	s.Send(s.outputSid, msg)
}

func (s *String) Connect(sid string) {
	s.outputSid = sid
}

func (s *String) Send(sid string, msg StringMessage) {
	//dispatcher_client.SendStringMessage(sid, msg) // debug code

	targetString := getString(sid)

	if targetString == nil { // string is not local, send msg to dispatcher
		dispatcher_client.SendStringMessage(sid, msg)
	} else { // found the target string on this vacuum server
		targetString.inputChan <- msg
	}
}

func (s *String) DeclareService(name string) {
	DeclareService(s.ID, name)
}

func (s *String) SendToService(serviceName string, msg StringMessage) {
	stringID := chooseServiceString(serviceName)
	s.Send(stringID, msg)
}
