package teststring

import "C"

import (
	"unsafe"

	"github.com/xiaonanln/vacuum/stringdev"
)

func TestString1(_s unsafe.Pointer) {
	s := *(*stringdev.StringInt)(_s)
	vlog.Debugf("TestString1 %s", s.ID())

	msg := s.Read()
	vlog.Debugln("Read", msg)
	s.Output(msg)
	return
}

func TestString2(s stringdev.StringInt) {
	vlog.Debugf("TestString2 %s", s.ID())

	msg := s.Read()
	vlog.Debugln("Read", msg)
	s.Output(msg)
	return
}
