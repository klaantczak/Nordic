package rnd

import (
	"hps"
	"math/rand"
	"time"
)

type randEx rand.Rand

func (r *randEx) Next() float64 {
	return (*rand.Rand)(r).Float64()
}

func Default() hps.IRnd {
	return hps.IRnd((*randEx)(rand.New(rand.NewSource(time.Now().UnixNano()))))
}
