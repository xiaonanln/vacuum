package main

import (
	log "github.com/Sirupsen/logrus"

	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	CALCULATOR_COUNT = 3
	RECEIVER_COUNT   = 1
	SENDER_COUNT     = 1
)

func Main(s *vacuum.String) {
	log.Debugf("Main %v running ...", s)
	s.DeclareService("Dispatcher") // declare the dispatcher service

	for i := 0; i < SENDER_COUNT; i++ {
		vacuum.CreateString("Sender")
	}
	for i := 0; i < RECEIVER_COUNT; i++ {
		vacuum.CreateString("Receiver")
	}
}

func Sender(s *vacuum.String) {
	log.Debugf("Sender running ...")
	s.DeclareService("Sender") // declare the dispatcher service

	// wait for receivers to be ready
	for {
		receiverCounter := vacuum.GetServiceProviderCount("Receiver")
		log.Debugln("receiverCounter", receiverCounter)
		if receiverCounter >= RECEIVER_COUNT {
			break
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Debugf("Now, all receivers are ready!!!")
	// wait until there are any senders
	for i := 0; i < 10; i++ {
		log.Debugln("Send", i)
		s.SendToService("Receiver", i)
	}
}

func Receiver(s *vacuum.String) {
	log.Debugf("Receiver running ...")
	s.DeclareService("Receiver") // declare the dispatcher service

	for {
		msg := s.Read()
		log.Debugln("Receive", msg)
	}
}

func main() {
	vacuum.RegisterString("Main", Main) // main string is a special string for system to boot up
	vacuum.RegisterString("Sender", Sender)
	vacuum.RegisterString("Receiver", Receiver)

	vacuum_server.RunServer()
}
