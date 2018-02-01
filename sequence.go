package main

import (
	"sync/atomic"
)

type Sequence struct {
	seq uint64
}

// NewSequence new a Sequence struct. store the seq.
func NewSequence() *Sequence {
	s := new(Sequence)
	s.seq = 0
	return s
}

func (s *Sequence) nextSeq() uint64 {
	atomic.AddUint64(&s.seq, 1);
	return s.seq
}