package msgbufpool

import "sync"

const (
	MSGBUF_SIZE   = 1 * 1024 * 1024
	MAX_POOL_SIZE = 100
)

type Msgbuf_t [MSGBUF_SIZE]byte

type msgbufpoolPerfData struct {
	getHit   uint64
	getMiss  uint64
	freeHit  uint64
	freeMiss uint64
}

var (
	pool     = make([]*Msgbuf_t, 0, MAX_POOL_SIZE)
	poolLock = sync.Mutex{}
	perfData msgbufpoolPerfData
)

func init() {

}

func GetMsgBuf() *Msgbuf_t {
	var msgbuf *Msgbuf_t
	poolLock.Lock()

	if len(pool) > 0 {
		perfData.getHit += 1

		last := len(pool) - 1
		msgbuf = pool[last]
		pool = pool[:last]

	} else {
		perfData.getMiss += 1
		msgbuf = &Msgbuf_t{} // allocate new one
	}

	poolLock.Unlock()
	return msgbuf
}

func PutMsgBuf(mb *Msgbuf_t) {
	poolLock.Lock()
	if len(pool) < MAX_POOL_SIZE {
		perfData.freeHit += 1
		pool = append(pool, mb)
	} else {
		// too many msgbuf in pool, we drop this one
		perfData.freeMiss += 1
	}
	poolLock.Unlock()
}
