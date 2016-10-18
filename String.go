package vacuum

import (
	"fmt"
	"log"

	. "github.com/xiaonanln/vacuum/common"

	"github.com/xiaonanln/vacuum/uuid"
)

type StringRoutine func(*String)

type String struct {
	ID        string
	routine   StringRoutine
	inputChan chan StringMessage
	outputSid string
}

func newString(routine StringRoutine) *String {
	return &String{
		ID:        uuid.GenUUID(),
		routine:   routine,
		inputChan: make(chan StringMessage),
		outputSid: "",
	}
}

func (s *String) String() string {
	return fmt.Sprintf("String<%s>@%v", s.ID, s.routine)
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
	targetString := getString(sid)

	if targetString != nil { // found the target string on this vacuum server
		targetString.inputChan <- msg
	} else {
		// output string not set, just write to output
		log.Printf("%s OUTPUT %T(%v)", s, msg, msg)
	}
}

func (s *String) DeclareService(name string) {
	DeclareService(s.ID, name)
}

func (s *String) SendToService(serviceName string, msg StringMessage) {
	stringID := chooseServiceString(serviceName)
	s.Send(stringID, msg)
}
