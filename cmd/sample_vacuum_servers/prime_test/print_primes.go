package main

import (
	"log"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

var (
	numberGeneratorID = ""
)

func Main(s *vacuum.String) {
	numberGeneratorID = vacuum.CreateString("NumberGenerator")
	log.Printf("NumBenerator created with ID: %s", numberGeneratorID)
}

func NumberGenerator(s *vacuum.String) {

}

func PrimeTester(s *vacuum.String) {

}

func PrimeOutputer(s *vacuum.String) {

}

func main() {
	vacuum.RegisterString("Main", Main)
	vacuum.RegisterString("NumberGenerator", NumberGenerator)
	vacuum.RegisterString("PrimeTester", PrimeTester)
	vacuum.RegisterString("PrimeOutputer", PrimeOutputer)
	vacuum_server.RunServer()
}
