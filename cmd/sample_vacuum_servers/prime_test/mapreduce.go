package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"
	"github.com/xiaonanln/vacuum/mapreduce"
	"github.com/xiaonanln/vacuum/vacuum_server"
	"gopkg.in/xiaonanln/typeconv.v0"
)

const (
	MAPPER_COUNT = 1
	BATCH_SIZE   = 10000
	BATCH_COUNT  = 3
)

func PrimeTestByMapReduceMain(s *vacuum.String) {
	for i := 0; i < MAPPER_COUNT; i++ {
		mapreduce.CreateMap("GetPrimesBetween", "CollectAllPrimes")
	}
	mapreduce.WaitReady("GetPrimesBetween", MAPPER_COUNT)

	mapreduce.CreateReduce("CollectAllPrimes", nil, "")
	mapreduce.WaitReady("CollectAllPrimes", 1)

	for i := 1; i < BATCH_COUNT; i++ {
		mapreduce.Send("GetPrimesBetween", []int{(i-1)*BATCH_SIZE + 1, (i) * BATCH_SIZE})
		s.Yield()
	}
	// wait for all mappers to quit
	mapreduce.Broadcast("GetPrimesBetween", nil)
	mapreduce.WaitGone("GetPrimesBetween")       // wait all GetPrimesBetween to finish
	mapreduce.Broadcast("CollectAllPrimes", nil) // send nil to CollectAllPrimes
	mapreduce.WaitGone("CollectAllPrimes")
}

func GetPrimesBetween(input interface{}) interface{} {
	range_ := typeconv.IntTuple(input)
	minNum := range_[0]
	maxNum := range_[1]

	logrus.Printf("GetPrimesBetween: %d ~ %d", minNum, maxNum)
	primes := []int64{}
	for i := minNum; i <= maxNum; i++ {
		if prime.IsPrime(int(i)) {
			primes = append(primes, int64(i))
		}
	}
	return primes
}

func CollectAllPrimes(_accum interface{}, _input interface{}) interface{} {
	var accum []int64
	if _accum != nil {
		accum = typeconv.IntTuple(_accum)
	}

	input := typeconv.IntTuple(_input)

	for _, n := range input {
		accum = append(accum, n)
	}

	logrus.Printf("Primes count: %d", len(accum))
	return accum
}

func main() {
	logrus.Debugf("Prime test usign map-reduce...")
	vacuum.RegisterMain(PrimeTestByMapReduceMain)
	mapreduce.RegisterMapFunc("GetPrimesBetween", GetPrimesBetween)
	mapreduce.RegisterReduceFunc("CollectAllPrimes", CollectAllPrimes)

	vacuum_server.RunServer()
}
