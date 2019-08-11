// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"math"
	"testing"

	. "github.com/flyingmutant/rapid"
)

// wrapper to test (*T).Helper()
func fatalf(t *T, format string, args ...interface{}) {
	t.Helper()
	t.Fatalf(format, args...)
}

func TestFailure_ImpossibleData(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		_ = Ints().Filter(func(i int) bool { return false }).Draw(t, "i")
	})
}

func TestFailure_Trivial(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		i := Ints().Draw(t, "i").(int)
		if i > 1000000000 {
			fatalf(t, "got a huge integer: %v", i)
		}
	})
}

func TestFailure_SimpleCollection(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		s := SlicesOf(Ints().Filter(func(i int) bool { return i%2 == -1 })).Draw(t, "s").([]int)
		if len(s) > 3 {
			fatalf(t, "got a long sequence: %v", s)
		}
	})
}

func TestFailure_CollectionElements(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		s := SlicesOfN(Ints(), 2, -1).Draw(t, "s").([]int)

		n := 0
		for _, i := range s {
			if i > 1000000 {
				n++
			}
		}

		if n > 1 {
			fatalf(t, "got %v huge elements", n)
		}
	})
}

func TestFailure_TrivialString(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		s := Strings().Draw(t, "s").(string)
		if len(s) > 7 {
			fatalf(t, "got bad string %v", s)
		}
	})
}

func TestFailure_Make(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		n := IntsMin(0).Draw(t, "n").(int)
		_ = make([]int, n)
	})
}

func TestFailure_Mean(t *testing.T) {
	t.Skip("expected failure")

	Check(t, func(t *T) {
		s := SlicesOf(Float64s()).Draw(t, "s").([]float64)

		mean := 0.0
		for _, f := range s {
			mean += f
		}
		mean /= float64(len(s))

		min, max := math.Inf(0), math.Inf(-1)
		for _, f := range s {
			if f < min {
				min = f
			}
			if f > max {
				max = f
			}
		}

		if mean < min || mean > max {
			t.Fatalf("got mean %v for range [%v, %v]", mean, min, max)
		}
	})
}

func TestFailure_ExampleParseDate(t *testing.T) {
	t.Skip("expected failure")

	Example_parseDate(t)
}

func TestFailure_ExampleQueue(t *testing.T) {
	t.Skip("expected failure")

	Example_queue(t)
}
