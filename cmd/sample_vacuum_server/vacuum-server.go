package main

import (
	"math/rand"
	"time"

	"log"

	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/vacuum_server"
)

const (
	CALCULATOR_COUNT = 3
)

func init() {
}

func dispatcher(_ *vacuum.String) {
	summer, _ := vacuum.CreateString(summer)
	log.Printf("Summer String created: %s", summer)

	calculators := []*vacuum.String{}

	for i := 0; i < CALCULATOR_COUNT; i++ {
		calculator, _ := vacuum.CreateString(calculator)
		log.Printf("Calculator String created: %s", calculator)
		calculators = append(calculators, calculator)
		calculator.Connect(summer.ID)
	}

	chooseRandomCalculator := func() *vacuum.String {
		i := rand.Intn(len(calculators))
		return calculators[i]
	}

	for i := 0; i < 10000; i++ {
		time.Sleep(10 * time.Millisecond)
		calculator := chooseRandomCalculator()
		calculator.Input(i)
	}
	log.Println("DISPATCH DONE")
}

func calculator(s *vacuum.String) {
	for {
		msg := s.Read()
		val := msg.(int)

		//for i := 0; i < 10000; i++ {
		//	val = val * val
		//}
		s.Output(val)
	}
}

func summer(s *vacuum.String) {
	log.Println("summer started!!!")
	var totalVal uint64 = 0
	var nextGrade uint64 = 100000

	for {
		msg := s.Read()

		val := msg.(int)
		totalVal += uint64(val)
		if totalVal >= nextGrade {
			s.Output(totalVal)
			nextGrade += 100000
		}
	}
}

func main() {
	s, _ := vacuum.CreateString(dispatcher)
	log.Printf("dispatcher created: %s", s)
	vacuum_server.RunServer()
}
