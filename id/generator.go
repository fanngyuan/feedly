package id

import (
	"sync/atomic"
)

type IdGenerable interface{
	GetId() uint64
}

type IdGenerator struct {
	IdSequence *int64
}

func (this IdGenerator) GetId() uint64 {
	seq:=atomic.AddInt64(this.IdSequence,1)
	return uint64(seq)
}
