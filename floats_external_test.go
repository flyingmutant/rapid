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
		Float32(),
		Float32Min(-0.1),
		Float32Min(1),
		Float32Max(0.1),
		Float32Max(2.5),
		Float32Range(0.3, 0.30001),
		Float32Range(0.3, 0.301),
		Float32Range(0.3, 0.7),
		Float32Range(math.E, math.Pi),
		Float32Range(0, 1),
		Float32Range(1, 2.5),
		Float32Range(0, 100),
		Float32Range(0, 10000),
		Float64(),
		Float64Min(-0.1),
		Float64Min(1),
		Float64Max(0.1),
		Float64Max(2.5),
		Float64Range(0.3, 0.30000001),
		Float64Range(0.3, 0.301),
		Float64Range(0.3, 0.7),
		Float64Range(math.E, math.Pi),
		Float64Range(0, 1),
		Float64Range(1, 2.5),
		Float64Range(0, 100),
		Float64Range(0, 10000),
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

	Check(t, func(t *T) {
		min := Float32().Draw(t, "min").(float32)
		max := Float32().Draw(t, "max").(float32)
		if min > max {
			min, max = max, min
		}

		g := Float32Range(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 400; i++ {
			f_, _, _ := g.Example(uint64(i))
			f := f_.(float32)

			gotMin = gotMin || f == min
			gotMax = gotMax || f == max
			gotZero = gotZero || f == 0

			if gotMin && gotMax && (min > 0 || max < 0 || gotZero) {
				return
			}
		}

		t.Fatalf("[%v, %v]: got min %v, got max %v, got zero %v", min, max, gotMin, gotMax, gotZero)
	})
}

func TestFloat64sBoundCoverage(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T) {
		min := Float64().Draw(t, "min").(float64)
		max := Float64().Draw(t, "max").(float64)
		if min > max {
			min, max = max, min
		}

		g := Float64Range(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 400; i++ {
			f_, _, _ := g.Example(uint64(i))
			f := f_.(float64)

			gotMin = gotMin || f == min
			gotMax = gotMax || f == max
			gotZero = gotZero || f == 0

			if gotMin && gotMax && (min > 0 || max < 0 || gotZero) {
				return
			}
		}

		t.Fatalf("[%v, %v]: got min %v, got max %v, got zero %v", min, max, gotMin, gotMax, gotZero)
	})
}
