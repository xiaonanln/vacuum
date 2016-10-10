package main

import (
	"math/rand"
	"time"

	"log"

	"github.com/xiaonanln/vacuum"
)

const (
	CALCULATOR_COUNT = 100
)

func dispatcher(s *vacuum.String) {
	summer, _ := vacuum.CreateString(summer)
	log.Printf("Summer String created: %s", summer)
	s.Input(100)

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
}

func calculator(s *vacuum.String) {
	for {
		msg := s.Read()
		log.Println("calculator!!!")
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
	for {
		msg := s.Read()
		log.Println("summer!!!")

		val := msg.(int)
		totalVal += uint64(val)
		s.Output(totalVal)
	}
}

func main() {
	s, _ := vacuum.CreateString(dispatcher)
	log.Printf("dispatcher created: %s", s)
	for {
		mainloop()
	}
}

func mainloop() {
	time.Sleep(time.Second)
}
