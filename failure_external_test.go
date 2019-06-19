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
	t.Skip()

	Check(t, func(t *T, i int) {}, Ints().Filter(func(i int) bool { return false }))
}

func TestFailure_Trivial(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, i int) {
		if i > 1000000000 {
			fatalf(t, "got a huge integer: %v", i)
		}
	}, Ints())
}

func TestFailure_SimpleCollection(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, s []int) {
		if len(s) > 3 {
			fatalf(t, "got a long sequence: %v", s)
		}
	}, SlicesOf(Ints().Filter(func(i int) bool { return i%2 == -1 })))
}

func TestFailure_CollectionElements(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, s []int) {
		n := 0
		for _, i := range s {
			if i > 1000000 {
				n++
			}
		}

		if n > 1 {
			fatalf(t, "got %v huge elements", n)
		}
	}, SlicesOfN(Ints(), 2, -1))
}

func TestFailure_TrivialString(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, s string) {
		if len(s) > 7 {
			fatalf(t, "got bad string %v", s)
		}
	}, Strings())
}

func TestFailure_Make(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, n int) {
		_ = make([]int, n)
	}, IntsMin(0))
}

func TestFailure_Mean(t *testing.T) {
	t.Skip()

	Check(t, func(t *T, s []float64) {
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
	}, SlicesOf(Float64s()))
}

func TestFailure_ExampleParseDate(t *testing.T) {
	t.Skip()

	Example_parseDate(t)
}

func TestFailure_ExampleQueue(t *testing.T) {
	t.Skip()

	Example_queue(t)
}
