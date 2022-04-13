package sequence

import "strconv"

type Seq uint

var GlobalSeq Seq = 0

func (seq *Seq) NextUInt() uint {
	*seq++
	return uint(*seq)
}

func (seq *Seq) NextInt() int {
	return int(seq.NextUInt())
}

func (seq *Seq) NextString() string {
	return strconv.Itoa(seq.NextInt())
}
