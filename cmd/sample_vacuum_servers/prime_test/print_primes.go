package main

import (
	log "github.com/Sirupsen/logrus"

	"math"

	"fmt"

	"time"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	PRIME_TESTER_COUNT = 10
	BATCH_SIZE         = 10000
)

var (
	startTime time.Time
	endTime   time.Time
)

func isPrimaryServer() bool {
	return vacuum_server.ServerID() == 1
}

func Main(s *vacuum.String) {
	if isPrimaryServer() {
		log.Printf("THIS IS THE PRIMARY SERVER")
		for i := 0; i < PRIME_TESTER_COUNT; i++ {
			vacuum.CreateString("PrimeTester")
		}
		vacuum.WaitServiceReady("PrimeTester", PRIME_TESTER_COUNT)

		vacuum.CreateString("BatchGenerator")
		vacuum.WaitServiceReady("BatchGenerator", 1) // all servers need to wait for BatchGenerator

		vacuum.CreateString("PrimeOutputer")
		vacuum.WaitServiceReady("PrimeOutputer", 1)
	} else {
		log.Printf("THIS IS SERVER %d", vacuum_server.ServerID())
		vacuum.WaitServiceReady("PrimeTester", PRIME_TESTER_COUNT)
		vacuum.WaitServiceReady("BatchGenerator", 1) // all servers need to wait for BatchGenerator
		vacuum.WaitServiceReady("PrimeOutputer", 1)
	}

}

func BatchGenerator(s *vacuum.String) {
	s.DeclareService("BatchGenerator")

	n := 1
	for {
		s.SendToService("PrimeTester", []int{
			n, n + BATCH_SIZE - 1,
		})
		n += BATCH_SIZE
	}
}

func PrimeTester(s *vacuum.String) {
	s.DeclareService("PrimeTester")

	for {
		n := s.Read().()
		if prime.IsPrime(n) {
			s.SendToService("PrimeOutputer", n)
		}
	}
}

func PrimeOutputer(s *vacuum.String) {
	s.DeclareService("PrimeOutputer")
	count := 0
	for {
		num := s.ReadInt()
		count += 1
		fmt.Printf("%d ", num)
		if count%20 == 0 {
			fmt.Print("\n")
		}
	}
}

func main() {
	//measureDirectCalculation()
	vacuum.RegisterString("Main", Main)
	vacuum.RegisterString("BatchGenerator", BatchGenerator)
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
