package gameserver

import "github.com/xiaonanln/vacuum/ext/entity"

const ()

type GSGate struct {
	entity.Entity
}

func init() {
	entity.RegisterEntity("GSGate", &GSGate{})
}

func runGates(num int) {
	for i := 0; i < num; i++ {
		entity.CreateEntity("GSGate")
	}
}

func (gate *GSGate) Init() {

}
