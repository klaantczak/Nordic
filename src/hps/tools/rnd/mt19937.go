// mt19937.go - an implementation of the 64bit Mersenne Twister PRNG
// Copyright (C) 2013  Jochen Voss <voss@seehuhn.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package rnd

const (
	n         = 624
	m         = 397
	notSeeded = n + 1

	hiMask uint32 = 0x80000000
	loMask uint32 = 0x7fffffff

	matrixA uint32 = 0x9908b0df
)

// MT19937 is the structure to hold the state of one instance of the
// Mersenne Twister PRNG.  New instances can be allocated using the
// mt19937.New() function.  MT19937 implements the rand.Source
// interface and rand.New() from the math/rand package can be used to
// generate different distributions from a MT19937 PRNG.
//
// This class is not safe for concurrent accesss by different
// goroutines.  If more than one goroutine accesses the PRNG, the
// callers must synchronise access using sync.Mutex or similar.
type MT19937 struct {
	state []uint32
	index int
}

// New allocates a new instance of the 64bit Mersenne Twister.
// A seed can be set using the .Seed() or .SeedFromSlice() methods.
func MT19937New(seed uint32) *MT19937 {
	res := &MT19937{
		state: make([]uint32, n),
		index: notSeeded,
	}
	res.Seed(seed)
	return res
}

// Seed uses the given 64bit value to initialise the generator state.
// This method is part of the rand.Source interface.
func (mt *MT19937) Seed(seed uint32) {
	x := mt.state
	x[0] = seed
	for i := uint32(1); i < n; i++ {
		x[i] = 1812433253*(x[i-1]^(x[i-1]>>30)) + i
	}
	mt.index = n
}

// Uint64 generates a (pseudo-)random 32bit value.  The output can be
// used as a replacement for a sequence of independent, uniformly
// distributed samples in the range 0, 1, ..., 2^32-1.
func (mt *MT19937) UInt32() uint32 {
	x := mt.state
	if mt.index >= n {
		if mt.index == notSeeded {
			mt.Seed(5489) // default seed, as in mt19937-64.c
		}
		for i := 0; i < n-m; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+m] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		for i := n - m; i < n-1; i++ {
			y := (x[i] & hiMask) | (x[i+1] & loMask)
			x[i] = x[i+(m-n)] ^ (y >> 1) ^ ((y & 1) * matrixA)
		}
		y := (x[n-1] & hiMask) | (x[0] & loMask)
		x[n-1] = x[m-1] ^ (y >> 1) ^ ((y & 1) * matrixA)
		mt.index = 0
	}
	y := x[mt.index]
	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18
	mt.index++
	return y
}

func (mt *MT19937) Float3() float64 {
	return (float64(mt.UInt32()) + 0.5) * (1.0 / 4294967296.0)
}
