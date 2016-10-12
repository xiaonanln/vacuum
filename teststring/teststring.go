package teststring

import (
	"log"

	"github.com/xiaonanln/vacuum/stringdev"
)

func TestString(s stringdev.StringInt) {
	log.Printf("TestString %s", s.ID())
	msg := s.Read()
	log.Println("Read", msg)
	s.Output(msg)
	return
}
