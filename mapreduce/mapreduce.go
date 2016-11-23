package mapreduce

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/xiaonanln/vacuum"
	"github.com/xiaonanln/vacuum/common"
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
	vacuum.RegisterString(MAPPER_STRING_NAME, makeMapper)
	vacuum.RegisterString(REDUCER_STRING_NAME, makeReducer)
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

type Mapper struct {
	funcName          string
	mapFunc           MapFunc
	outputServiceName string
}

func (m *Mapper) Init(s *vacuum.String) {
	m.funcName = s.Args()[0].(string)
	m.outputServiceName = getServiceName(s.Args()[1].(string))
	m.mapFunc = mapFuncs[m.funcName]
	s.DeclareService(getServiceName(m.funcName))
}

func (m *Mapper) Fini(s *vacuum.String) {}

func (m *Mapper) Loop(s *vacuum.String, input common.StringMessage) {
	output := m.mapFunc(input)
	// send the output to the next Mapper / Reducer
	if m.outputServiceName != "" {
		s.SendToService(m.outputServiceName, output)
	} else {
		logrus.Printf("Mapper %s output: %v", m.funcName, output)
	}
}

func makeMapper() vacuum.StringDelegate {
	return &Mapper{}
}

type Reducer struct {
	funcName          string
	outputServiceName string
	accum             interface{}
	reduceFunc        ReduceFunc
}

func (r *Reducer) Init(s *vacuum.String) {
	args := s.Args()
	r.funcName = args[0].(string)
	r.outputServiceName = getServiceName(args[1].(string))
	r.accum = args[2]
	r.reduceFunc = reduceFuncs[r.funcName]

	s.DeclareService(getServiceName(r.funcName))
}

func (r *Reducer) Fini(s *vacuum.String) {
	if r.outputServiceName != "" {
		s.SendToService(r.outputServiceName, r.accum)
	} else {
		fmt.Printf("%s: %v\n", r.funcName, r.accum)
	}
}

func (r *Reducer) Loop(s *vacuum.String, msg common.StringMessage) {
	r.accum = r.reduceFunc(r.accum, msg)
}

func makeReducer() vacuum.StringDelegate {
	return &Reducer{}
}

func CreateMap(funcName string, outputFuncName string) {
	vacuum.CreateString(MAPPER_STRING_NAME, funcName, outputFuncName)
}

func CreateReduce(funcName string, initial interface{}, outputFuncName string) {
	vacuum.CreateString(REDUCER_STRING_NAME, funcName, outputFuncName, initial)
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
