package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	PRIME_TESTER_COUNT = 16
	BATCH_SIZE         = 10000
)

func isPrimaryServer() bool {
	return vacuum_server.ServerID() == 1
}

func Main(s *vacuum.String) {
	if isPrimaryServer() {
		log.Infof("THIS IS THE PRIMARY SERVER")
		vacuum.CreateString("PrimeOutputer")
		vacuum.WaitServiceReady("PrimeOutputer", 1)

		for i := 0; i < PRIME_TESTER_COUNT; i++ {
			vacuum.CreateString("PrimeTester")
		}
		vacuum.WaitServiceReady("PrimeTester", PRIME_TESTER_COUNT)

		vacuum.CreateString("BatchGenerator")
		vacuum.WaitServiceReady("BatchGenerator", 1) // all servers need to wait for BatchGenerator

	} else {
		log.Infof("THIS IS SERVER %d", vacuum_server.ServerID())
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
		primes := []int{}
		range_ := s.ReadIntTuple()
		//log.Debugf("PrimeTester: testing %d ~ %d ...", range_[0], range_[1])
		for n := range_[0]; n <= range_[1]; n++ {
			if prime.IsPrime(n) {
				primes = append(primes, n)
			}
		}
		s.SendToService("PrimeOutputer", primes)
	}
}

type _SortedOutput [][]int

//
//func (L _SortedOutput) Len() int {
//	return len(L)
//}
//
//func (L _SortedOutput) Less(i, j int) bool {
//	return L[i][0] < L[j][0]
//}

func (L _SortedOutput) Swap(i, j int) {
	var tmp []int
	tmp = L[i]
	L[i] = L[j]
	L[j] = tmp
}

func PrimeOutputer(s *vacuum.String) {
	s.DeclareService("PrimeOutputer")
	//expectNum := 1
	//sortedOutputs := _SortedOutput{}

	for {
		nums := s.ReadIntTuple()
		log.Debugf("PrimeOutputer: Read %v", nums)
		//sortedOutputs = append(sortedOutputs, nums)
		//sort.Sort(sortedOutputs)
		//for len(sortedOutputs) > 0 && expectNum == sortedOutputs[0][0] {
		//	for _, n := range
		//	expectNum = sortedOutputs[0][1] + 1
		//	sortedOutputs = sortedOutputs[1:]
		//}

		for _, n := range nums {
			prime.OutputPrime(n)
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
