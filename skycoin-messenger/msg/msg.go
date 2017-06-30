package msg

import (
	"sync"
	"github.com/skycoin/net/client"
)

var (
	OP_POOL       = make([]*sync.Pool, OP_SIZE)
	op_pool_mutex = new(sync.RWMutex)
)

type OP interface {
	Execute(OPer) error
}

type OPer interface {
	GetFactory() *client.ClientConnectionFactory
	SetFactory(factory *client.ClientConnectionFactory)
	PushLoop(conn *client.ClientConnection, data []byte)
}

func GetOP(opn int) OP {
	if opn < 0 || opn > OP_SIZE {
		return nil
	}

	op_pool_mutex.RLock()
	op, ok := OP_POOL[opn].Get().(OP)
	op_pool_mutex.RUnlock()
	if !ok {
		return nil
	}
	return op
}

func PutOP(opn int, op OP) {
	if opn < 0 || opn > OP_SIZE {
		return
	}
	op_pool_mutex.Lock()
	OP_POOL[opn].Put(op)
	op_pool_mutex.Unlock()
}