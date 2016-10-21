package teststring

import "C"

import (
	"unsafe"

	log "github.com/Sirupsen/logrus"

	"github.com/xiaonanln/vacuum/stringdev"
)

func TestString1(_s unsafe.Pointer) {
	s := *(*stringdev.StringInt)(_s)
	log.Debugf("TestString1 %s", s.ID())

	msg := s.Read()
	log.Debugln("Read", msg)
	s.Output(msg)
	return
}

func TestString2(s stringdev.StringInt) {
	log.Debugf("TestString2 %s", s.ID())

	msg := s.Read()
	log.Debugln("Read", msg)
	s.Output(msg)
	return
}
