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
	log.Printf("%s INPUT %T(%v)", s, msg, msg)
	s.inputChan <- msg
}

func (s *String) Read() StringMessage {
	msg := <-s.inputChan
	log.Printf("%s READ %T(%v)", s, msg, msg)
	return msg
}

func (s *String) Output(msg StringMessage) {
	outputString := getString(s.outputSid)
	log.Printf("get output string: %s", outputString)

	if outputString != nil {
		outputString.Input(msg)
	} else {
		// output string not set, just write to output
		fmt.Printf("%s OUTPUT %T(%v)\n", s, msg, msg)
	}
}

func (s *String) Connect(sid string) {
	s.outputSid = sid
}
