package main

import (
	"log"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	CALCULATOR_COUNT = 3
)

var ()

func init() {
}

func Sender(s *vacuum.String) {
	log.Printf("Sender running ...")
}

func Receiver(s *vacuum.String) {
	log.Printf("Receiver running ...")
}

func main() {
	vacuum.RegisterString("Sender", Sender)
	vacuum.RegisterString("Receiver", Receiver)

	vacuum.CreateString("Sender")
	vacuum.CreateString("Receiver")

	vacuum_server.RunServer()
}
