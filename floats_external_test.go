// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"math"
	"sort"
	"testing"

	. "github.com/flyingmutant/rapid"
)

func TestFloatsExamples(t *testing.T) {
	gens := []*Generator{
		Float32s(),
		Float32sMin(-0.1),
		Float32sMin(1),
		Float32sMax(0.1),
		Float32sMax(2.5),
		Float32sRange(0.3, 0.30001),
		Float32sRange(0.3, 0.301),
		Float32sRange(0.3, 0.7),
		Float32sRange(math.E, math.Pi),
		Float32sRange(0, 1),
		Float32sRange(1, 2.5),
		Float32sRange(0, 100),
		Float32sRange(0, 10000),
		Float64s(),
		Float64sMin(-0.1),
		Float64sMin(1),
		Float64sMax(0.1),
		Float64sMax(2.5),
		Float64sRange(0.3, 0.30000001),
		Float64sRange(0.3, 0.301),
		Float64sRange(0.3, 0.7),
		Float64sRange(math.E, math.Pi),
		Float64sRange(0, 1),
		Float64sRange(1, 2.5),
		Float64sRange(0, 100),
		Float64sRange(0, 10000),
	}

	for _, g := range gens {
		t.Run(g.String(), func(t *testing.T) {
			var vals []float64
			var vals32 bool
			for i := 0; i < 100; i++ {
				f, _, _ := g.Example()
				_, vals32 = f.(float32)
				vals = append(vals, rv(f).Float())
			}
			sort.Float64s(vals)

			for _, f := range vals {
				if vals32 {
					t.Logf("%30g %10.3g % 5d % 20d % 16x", f, f, int(math.Log10(math.Abs(f))), int64(f), math.Float32bits(float32(f)))
				} else {
					t.Logf("%30g %10.3g % 5d % 20d % 16x", f, f, int(math.Log10(math.Abs(f))), int64(f), math.Float64bits(f))
				}
			}
		})
	}
}

func TestFloat32sBoundCoverage(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T, min float32, max float32) {
		if min > max {
			t.Skip("min > max")
		}

		g := Float32sRange(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 400; i++ {
			f_, _, _ := g.Example(uint64(i))
			f := f_.(float32)

			if f == min {
				gotMin = true
			}
			if f == max {
				gotMax = true
			}
			if f == 0 {
				gotZero = true
			}
			if gotMin && gotMax && (min > 0 || max < 0 || gotZero) {
				return
			}
		}

		t.Fatalf("[%v, %v]: got min %v, got max %v, got zero %v", min, max, gotMin, gotMax, gotZero)
	}, Float32s(), Float32s())
}

func TestFloat64sBoundCoverage(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T, min float64, max float64) {
		if min > max {
			t.Skip("min > max")
		}

		g := Float64sRange(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 400; i++ {
			f_, _, _ := g.Example(uint64(i))
			f := f_.(float64)

			if f == min {
				gotMin = true
			}
			if f == max {
				gotMax = true
			}
			if f == 0 {
				gotZero = true
			}
			if gotMin && gotMax && (min > 0 || max < 0 || gotZero) {
				return
			}
		}

		t.Fatalf("[%v, %v]: got min %v, got max %v, got zero %v", min, max, gotMin, gotMax, gotZero)
	}, Float64s(), Float64s())
}
