package simhashlsh

import (
	"testing"
)

func TestHyperplane(t *testing.T) {
	hs := newHyperplanes(2, 10)
	if len(hs) != 2 {
		t.Fail()
	}
	if len(hs[0]) != 10 {
		t.Fail()
	}
}
