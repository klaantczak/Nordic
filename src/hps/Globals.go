package hps

import (
	"math/rand"
)

var rnd IRnd

func SetRnd(r IRnd) {
	rnd = r
}

func GetRnd() IRnd {
	return rnd
}

func RndNext() float64 {
	if rnd == nil {
		return rand.Float64()
	}
	return rnd.Next()
}
