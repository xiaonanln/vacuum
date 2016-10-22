package main

import (
	"typeconv"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/cmd/sample_vacuum_servers/prime_test/internal/prime"
	"github.com/xiaonanln/vacuum/mapreduce"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

func Main(s *vacuum.String) {
	mapreduce.CreateMap("GetPrimesBetween", "")
	mapreduce.WaitMapper("GetPrimesBetween", 1)

	mapreduce.Send("GetPrimesBetween", []int{1, 10000})
}

func GetPrimesBetween(input interface{}) interface{} {
	range_ := typeconv.ToIntTuple(input)
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

func main() {
	logrus.Debugf("Prime test usign map-reduce...")
	vacuum.RegisterString("Main", Main)
	mapreduce.RegisterMapFunc("GetPrimesBetween", GetPrimesBetween)
	vacuum_server.RunServer()
}
