package sequence

import (
	"strconv"
	"sync/atomic"
)

type Seq struct {
	uint32
}

var GlobalSeq Seq = Seq{}

func (this *Seq) NextUint32() uint32 {
	return atomic.AddUint32(&this.uint32, 1)
}

func (this *Seq) NextUint() uint {
	return uint(this.NextUint32())
}

func (seq *Seq) NextInt() int {
	return int(seq.NextUint32())
}

func (seq *Seq) NextString() string {
	return strconv.Itoa(seq.NextInt())
}
