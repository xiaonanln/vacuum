package test

import (
	"testing"

	"sync"
)

type Message [1024 * 1024]byte

var (
	syncPool = sync.Pool{
		New: func() interface{} {
			return &Message{}
		},
	}
)

func init() {

}

func BenchmarkNewobject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcWithNewobject()
	}
}

func BenchmarkNoNewobject(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcWithoutNewobject()
	}
}

//func BenchmarkGetMsgBuf(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		funcUsingMsgbufpool()
//	}
//}

func funcWithNewobject() *Message {
	var b Message
	return &b
}

func funcWithoutNewobject() int {
	var b Message
	b[0] = 1
	return 1
}

//
//func funcUsingMsgbufpool() *Message {
//	t := msgbufpool.GetMsgBuf()
//	msgbufpool.PutMsgBuf(t)
//	return t
//}

func BenchmarkGetFromSyncPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcUsingSyncPool()
	}
}

func funcUsingSyncPool() {
	x := syncPool.Get()
	syncPool.Put(x)
}
