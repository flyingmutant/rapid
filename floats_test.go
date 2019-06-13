// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"testing"
)

func TestUfloat32FromParts(t *testing.T) {
	Check(t, func(t *T, f float32) {
		e, si, sf := ufloatParts(float64(f), float32ExpBits, float32SignifBits)
		g := float32(ufloatFromParts(float32SignifBits, e, si, sf))
		if g != f {
			t.Fatalf("got %v (0x%x) back from %v (0x%x)", g, math.Float32bits(g), f, math.Float32bits(g))
		}
	}, Float32sMin(0))
}

func TestUfloat64FromParts(t *testing.T) {
	Check(t, func(t *T, f float64) {
		e, si, sf := ufloatParts(f, float64ExpBits, float64SignifBits)
		g := ufloatFromParts(float64SignifBits, e, si, sf)
		if g != f {
			t.Fatalf("got %v (0x%x) back from %v (0x%x)", g, math.Float64bits(g), f, math.Float64bits(g))
		}
	}, Float64sMin(0))
}

func TestGenUfloat32Range(t *testing.T) {
	Check(t, func(t *T, min_ float32, max_ float32) {
		min := float64(min_)
		max := float64(max_)
		Assume(min <= max)
		f := genUfloatRange(t.src.s, min, max, float32ExpBits, float32SignifBits)
		if float64(float32(f)) != f {
			t.Fatalf("%v (0x%x) is not a float32", f, math.Float64bits(f))
		}
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float32sMin(0), Float32sMin(0))
}

func TestGenUfloat64Range(t *testing.T) {
	Check(t, func(t *T, min float64, max float64) {
		Assume(min <= max)
		f := genUfloatRange(t.src.s, min, max, float64ExpBits, float64SignifBits)
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64sMin(0), Float64sMin(0))
}

func TestGenFloat32Range(t *testing.T) {
	Check(t, func(t *T, min_ float32, max_ float32) {
		min := float64(min_)
		max := float64(max_)
		Assume(min <= max)
		f := genFloatRange(t.src.s, min, max, float32ExpBits, float32SignifBits)
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
		f := genFloatRange(t.src.s, min, max, float64ExpBits, float64SignifBits)
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64s(), Float64s())
}
