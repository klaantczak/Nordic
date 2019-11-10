package hpstest

import (
	"math/rand"
)

type TestRnd struct{}

func (t *TestRnd) Next() float64 {
	return rand.Float64()
}
