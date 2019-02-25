// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"testing"
)

func TestGenUfloat32Range(t *testing.T) {
	Check(t, func(t *T, min_ float32, max_ float32) {
		min := math.Abs(float64(min_))
		max := math.Abs(float64(max_))
		Assume(min <= max)
		f := genUfloatRange(t.data.s, min, max, float32ExpBits, float32SignifBits)
		if float64(float32(f)) != f {
			t.Fatalf("%v (0x%x) is not a float32", f, math.Float64bits(f))
		}
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float32s(), Float32s())
}

func TestGenUfloat64Range(t *testing.T) {
	Check(t, func(t *T, min float64, max float64) {
		min = math.Abs(min)
		max = math.Abs(max)
		Assume(min <= max)
		f := genUfloatRange(t.data.s, min, max, float64ExpBits, float64SignifBits)
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64s(), Float64s())
}

func TestGenFloat32Range(t *testing.T) {
	Check(t, func(t *T, min_ float32, max_ float32) {
		min := float64(min_)
		max := float64(max_)
		Assume(min <= max)
		f := genFloatRange(t.data.s, min, max, float32ExpBits, float32SignifBits)
		if float64(float32(f)) != f {
			t.Fatalf("%v (0x%x) is not a float32", f, math.Float64bits(f))
		}
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float32s(), Float32s())
}

func TestGenFloat64Range(t *testing.T) {
	Check(t, func(t *T, min float64, max float64) {
		Assume(min <= max)
		f := genFloatRange(t.data.s, min, max, float64ExpBits, float64SignifBits)
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64s(), Float64s())
}
