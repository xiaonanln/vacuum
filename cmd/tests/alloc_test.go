package test

import (
	"testing"

	"sync"

	"github.com/xiaonanln/vacuum/msgbufpool"
)

var (
	syncPool = sync.Pool{
		New: func() interface{} {
			return &msgbufpool.Msgbuf_t{}
		},
	}
)

func init() {
	msgbufpool.PutMsgBuf(msgbufpool.GetMsgBuf())
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

func BenchmarkGetMsgBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcUsingMsgbufpool()
	}
}

func funcWithNewobject() *msgbufpool.Msgbuf_t {
	var b msgbufpool.Msgbuf_t
	return &b
}

func funcWithoutNewobject() int {
	var b msgbufpool.Msgbuf_t
	b[0] = 1
	return 1
}

func funcUsingMsgbufpool() *msgbufpool.Msgbuf_t {
	t := msgbufpool.GetMsgBuf()
	msgbufpool.PutMsgBuf(t)
	return t
}

func BenchmarkGetFromSyncPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcUsingSyncPool()
	}
}

func funcUsingSyncPool() {
	x := syncPool.Get()
	syncPool.Put(x)
}
