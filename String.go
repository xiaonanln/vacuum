package vacuum

import (
	"fmt"

	"log"

	. "github.com/xiaonanln/vacuum/common"
	"github.com/xiaonanln/vacuum/vacuum_server/dispatcher_client"
)

const (
	STRING_MESSAGE_BUFFER_SIZE = 0
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

func (s *String) Connect(stringID string) {
	s.outputSid = stringID
}

func (s *String) Send(stringID string, msg StringMessage) {
	//dispatcher_client.SendStringMessage(sid, msg) // debug code
	if stringID == "" {
		log.Panicf("%s.Send: stringID is empty", s)
	}

	targetString := getString(stringID)

	if targetString == nil { // string is not local, send msg to dispatcher
		dispatcher_client.SendStringMessage(stringID, msg)
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
