package main

import (
	"log"

	"github.com/xiaonanln/vacuum/vacuum"
)

func main() {
	s := vacuum.String{}
	log.Printf("String %T %v", s, s)
}
