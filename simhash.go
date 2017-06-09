package simhashlsh

import (
	"math/rand"
)

type signature []uint8

// Key is a way to index into a table.
type hashTableKey []uint8

// Value is an index into the input dataset.
type hashTableBucket []string

type hashTable map[uint64]hashTableBucket

// Represents a simhash signature - an array of hash values
type simhash struct {
	sig signature
}

// signature generates the simhash of an attribute (a bag of values) using the hyperplanes e
func newSimhash(hs hyperplanes, e []float64) *simhash {
	sig := newSignature(hs, e)
	return &simhash{
		sig: sig,
	}
}

// newSignature computes the signature for a simhash of input float array
func newSignature(hyperplanes hyperplanes, e []float64) signature {
	sigarr := make([]uint8, len(hyperplanes))
	for hix, h := range hyperplanes {
		var dp float64
		for k, v := range e {
			dp += h[k] * float64(v)
		}
		if dp >= 0 {
			sigarr[hix] = uint8(1)
		} else {
			sigarr[hix] = uint8(0)
		}
	}
	return sigarr
}

type hyperplanes [][]float64

//NewHyperplanes generates and initializes a set of d hyperplanes with s dimensions.
func newHyperplanes(d, s int) hyperplanes {
	hs := make([][]float64, d)
	for i := 0; i < d; i++ {
		v := make([]float64, s)
		for i := 0; i < s; i++ {
			n := rand.NormFloat64()
			v[i] = n
		}
		hs[i] = v
	}
	return hs
}
