package vacuum

import (
	"fmt"
	"log"

	"github.com/xiaonanln/vacuum/uuid"
)

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

func (s *String) run() {
	go s.routine(s)
}

func (s *String) String() string {
	return fmt.Sprintf("String<%s>@%v", s.ID, s.routine)
}

func (s *String) Input(msg StringMessage) {
	//log.Printf("%s INPUT %T(%v)", s, msg, msg)
	s.inputChan <- msg
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
	if targetString != nil {
		targetString.Input(msg)
	} else {
		// output string not set, just write to output
		log.Printf("%s OUTPUT %T(%v)", s, msg, msg)
	}
}
