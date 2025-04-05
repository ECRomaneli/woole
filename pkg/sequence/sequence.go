package sequence

import (
	"strconv"
	"sync/atomic"
)

type Seq struct {
	uint32
}

var GlobalSeq Seq = Seq{}

func (seq *Seq) NextUint32() uint32 {
	return atomic.AddUint32(&seq.uint32, 1)
}

func (seq *Seq) NextUint() uint {
	return uint(seq.NextUint32())
}

func (seq *Seq) NextInt() int {
	return int(seq.NextUint32())
}

func (seq *Seq) NextString() string {
	return strconv.Itoa(seq.NextInt())
}

func (seq *Seq) LastUint32() uint32 {
	return seq.uint32
}

func (seq *Seq) LastUint() uint {
	return uint(seq.LastUint32())
}

func (seq *Seq) LastInt() int {
	return int(seq.LastUint32())
}

func (seq *Seq) LastString() string {
	return strconv.Itoa(seq.LastInt())
}
