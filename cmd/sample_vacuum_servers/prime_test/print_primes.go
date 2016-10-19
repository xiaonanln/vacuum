package main

import (
	"log"

	"math"

	"fmt"

	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	PRIME_TESTER_COUNT = 100
	MIN_NUMBER         = 1000000
	MAX_NUMBER         = MIN_NUMBER + 1000000
)

var (
	numberGeneratorID = ""
	primeOutputerID   = ""
	startTime         time.Time
	endTime           time.Time
)

//
//func measureDirectCalculation() {
//	t0 := time.Now()
//	for n := MIN_NUMBER; n <= MAX_NUMBER; n++ {
//		if isPrime(n) {
//			//fmt.Println(n)
//		}
//	}
//	t1 := time.Now()
//	log.Printf("Direct calculation takes %v", t1.Sub(t0))
//}

func isPrimaryServer() bool {
	return vacuum_server.ServerID() == 1
}

func Main(s *vacuum.String) {
	if isPrimaryServer() {
		log.Printf("THIS IS THE PRIMARY SERVER")
		numberGeneratorID = vacuum.CreateString("NumberGenerator")
		primeOutputerID = vacuum.CreateString("PrimeOutputer")
		log.Printf("NumBenerator: %s, PrimeOutputer: %s", numberGeneratorID, primeOutputerID)

		for i := 0; i < PRIME_TESTER_COUNT; i++ {
			vacuum.CreateString("PrimeTester")
		}
	} else {
		log.Printf("THIS IS SERVER %d", vacuum_server.ServerID())
	}

	vacuum.WaitServiceReady("NumberGenerator", 1) // all servers need to wait for NumberGenerator
	vacuum.WaitServiceReady("PrimeTester", PRIME_TESTER_COUNT)
	vacuum.WaitServiceReady("PrimeOutputer", 1)

	if isPrimaryServer() {
		s.Send(numberGeneratorID, MIN_NUMBER)
		s.Send(numberGeneratorID, MAX_NUMBER)
	}

}

func NumberGenerator(s *vacuum.String) {
	s.DeclareService("NumberGenerator")

	minNum := s.ReadInt()
	maxNum := s.ReadInt()
	log.Printf("NumberGenerator: %d ~ %d", minNum, maxNum)

	for n := minNum; n <= maxNum; n++ {
		s.SendToService("PrimeTester", n)
	}
}

func PrimeTester(s *vacuum.String) {
	s.DeclareService("PrimeTester")
	for {
		n := s.ReadInt()
		if MIN_NUMBER == n {
			startTime = time.Now()
		} else if MAX_NUMBER == n {
			endTime = time.Now()
			log.Printf("Distributed strings takes: %v", (endTime.Sub(startTime)))
		}

		if isPrime(n) {
			s.SendToService("PrimeOutputer", n)
		}
	}
}

func PrimeOutputer(s *vacuum.String) {
	s.DeclareService("PrimeOutputer")

	for {
		num := s.ReadInt()
		fmt.Println(num)
	}
}

func main() {
	//measureDirectCalculation()

	vacuum.RegisterString("Main", Main)
	vacuum.RegisterString("NumberGenerator", NumberGenerator)
	vacuum.RegisterString("PrimeTester", PrimeTester)
	vacuum.RegisterString("PrimeOutputer", PrimeOutputer)
	vacuum_server.RunServer()
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 {
		return true
	}
	sqrt := int(math.Sqrt(float64(n)) + 0.000001)
	for i := 2; i <= sqrt; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
