package mapreduce

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
)

const (
	MAPPER_STRING_NAME  = "mapreduce.mapper"
	REDUCER_STRING_NAME = "mapreduce.reducer"
	SERVICE_NAME_PREFIX = "mapreduce.service."
)

var (
	mapFuncs    = map[string]MapFunc{}
	reduceFuncs = map[string]ReduceFunc{}
)

func init() {
	vacuum.RegisterString(MAPPER_STRING_NAME, mapperRoutine)
	vacuum.RegisterString(REDUCER_STRING_NAME, reducerRoutine)
}

type MapFunc func(input interface{}) interface{}
type ReduceFunc func(accum interface{}, input interface{}) interface{}

func RegisterMapFunc(name string, f MapFunc) {
	if _, ok := mapFuncs[name]; ok {
		logrus.Panicf("Mapper func %s already registered", name)
	}
	if _, ok := reduceFuncs[name]; ok {
		logrus.Panicf("%s is already registered as Reducer func", name)
	}

	mapFuncs[name] = f
}

func RegisterReduceFunc(name string, f ReduceFunc) {
	if _, ok := mapFuncs[name]; ok {
		logrus.Panicf("%s is already registered as Mapper func", name)
	}
	if _, ok := reduceFuncs[name]; ok {
		logrus.Panicf("Reducer func %s is already registered", name)
	}

	reduceFuncs[name] = f
}

func mapperRoutine(s *vacuum.String) {
	funcName := s.ReadString()
	outputFuncName := s.ReadString()

	myServiceName := getServiceName(funcName)
	s.DeclareService(myServiceName) // declare the service of this map
	outputServiceName := getServiceName(outputFuncName)
	mapFunc := mapFuncs[funcName]

	for {
		input := s.Read() // read input, whatever it is
		if input == nil {
			// nil means end of execution
			break
		}

		output := mapFunc(input)
		// send the output to the next Mapper / Reducer
		if outputServiceName != "" {
			s.SendToService(outputServiceName, output)
		} else {
			logrus.Printf("Mapper %s output: %v", funcName, output)
		}
	}
}

func reducerRoutine(s *vacuum.String) {
	funcName := s.ReadString()
	initial := s.Read()
	outputFuncName := s.ReadString()

	myServiceName := getServiceName(funcName)
	s.DeclareService(myServiceName) // declare the service of this map
	reduceFunc := reduceFuncs[funcName]
	outputServiceName := getServiceName(outputFuncName)

	accum := initial

	for {
		input := s.Read() // read input, whatever it is
		if input == nil {
			break
		}

		accum = reduceFunc(accum, input)
		// send the output to the next Mapper / Reducer
	}

	if outputServiceName != "" {
		s.SendToService(outputServiceName, accum)
	} else {
		fmt.Printf("%s: %v\n", funcName, accum)
	}
}

func CreateMap(funcName string, outputFuncName string) {
	mapperID := vacuum.CreateString(MAPPER_STRING_NAME)
	vacuum.Send(mapperID, funcName)       // send the mapper name
	vacuum.Send(mapperID, outputFuncName) // send the name of next func (can be a mapper or reducer)
}

func CreateReduce(funcName string, initial interface{}, outputFuncName string) {
	reducerID := vacuum.CreateString(REDUCER_STRING_NAME)
	vacuum.Send(reducerID, funcName)
	vacuum.Send(reducerID, initial)
	vacuum.Send(reducerID, outputFuncName)
}

// get the service name for Mappers / Reducers
func getServiceName(mapperOrReducerName string) string {
	if mapperOrReducerName != "" {
		return SERVICE_NAME_PREFIX + mapperOrReducerName
	} else {
		return ""
	}
}

func WaitReady(name string, n int) {
	vacuum.WaitServiceReady(getServiceName(name), n)
}

func WaitGone(name string) {
	vacuum.WaitServiceGone(getServiceName(name))
}

func Send(name string, val interface{}) {
	serviceName := getServiceName(name)
	vacuum.SendToService(serviceName, val)
}

func Broadcast(name string, val interface{}) {
	serviceName := getServiceName(name)
	vacuum.BroadcastToService(serviceName, val)
}

func GetCount(name string) int {
	return vacuum.GetServiceProviderCount(getServiceName(name))
}
