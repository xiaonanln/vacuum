package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/typeconv"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"
	"github.com/xiaonanln/vacuum/mapreduce"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

func Main(s *vacuum.String) {
	mapreduce.CreateMap("GetPrimesBetween", "PrintAllPrimes")
	mapreduce.WaitMapper("GetPrimesBetween", 1)
	mapreduce.CreateReduce("PrintAllPrimes", nil, "")
	mapreduce.WaitReducer("PrintAllPrimes", 1)

	mapreduce.Send("GetPrimesBetween", []int{1, 10000})
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

func PrintAllPrimes(accum interface{}, input interface{}) interface{} {
	primes := typeconv.IntTuple(input)

	for _, n := range primes {
		logrus.Printf("Prime %d", n)
	}
	return accum
}

func main() {
	logrus.Debugf("Prime test usign map-reduce...")
	vacuum.RegisterString("Main", Main)
	mapreduce.RegisterMapFunc("GetPrimesBetween", GetPrimesBetween)
	mapreduce.RegisterReduceFunc("PrintAllPrimes", PrintAllPrimes)

	vacuum_server.RunServer()
}
