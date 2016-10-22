package mapreduce

import (
	"github.com/xiaonanln/vacuum"
	. "github.com/xiaonanln/vacuum/common"
)

const (
	MAPPER_NAME  = "mapreduce.MAPPER"
	REDUCER_NAME = "mapreduce.REDUCER"
)

var (
	mapStrings    = StringSet{}
	reduceStrings = StringSet{}
)

func init() {
	vacuum.RegisterString(MAPPER_NAME, mapperRoutine)
	vacuum.RegisterString(REDUCER_NAME, reducerRoutine)
}

type MapFunc func(input interface{}) interface{}
type ReduceFunc func(accum interface{}, input interface{}) interface{}

func RegisterMapFunc(name string, f MapFunc) {
}

func RegisterReduceFunc(name string, f ReduceFunc) {
}

func mapperRoutine(s *vacuum.String) {

}

func reducerRoutine(s *vacuum.String) {

}

func Map(name string) {
	mapperID := vacuum.CreateString(MAPPER_NAME)
}
