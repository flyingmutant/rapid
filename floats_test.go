// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"testing"
)

func TestFloatConversionRoundtrip(t *testing.T) {
	Check(t, func(t *T) {
		u := uint32(t.src.s.drawBits(32))
		f := math.Float32frombits(u)
		Assume(!math.IsNaN(float64(f))) // we can get NaNs with different bit patterns back
		g := float32(float64(f))
		if g != f {
			t.Fatalf("got %v (0x%x) back from %v (0x%x)", g, math.Float32bits(g), f, math.Float32bits(f))
		}
	})
}

func TestUfloat32FromParts(t *testing.T) {
	Check(t, func(t *T, f float32) {
		g := ufloat32FromParts(ufloat32Parts(f))
		if g != f {
			t.Fatalf("got %v (0x%x) back from %v (0x%x)", g, math.Float32bits(g), f, math.Float32bits(f))
		}
	}, Float32sMin(0))
}

func TestUfloat64FromParts(t *testing.T) {
	Check(t, func(t *T, f float64) {
		g := ufloat64FromParts(ufloat64Parts(f))
		if g != f {
			t.Fatalf("got %v (0x%x) back from %v (0x%x)", g, math.Float64bits(g), f, math.Float64bits(f))
		}
	}, Float64sMin(0))
}

func TestGenUfloat32Range(t *testing.T) {
	Check(t, func(t *T, min float32, max float32) {
		Assume(min <= max)
		f := ufloat32FromParts(genUfloatRange(t.src.s, float64(min), float64(max), float32SignifBits))
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float32bits(f), min, max, math.Float32bits(min), math.Float32bits(max))
		}
	}, Float32sMin(0), Float32sMin(0))
}

func TestGenUfloat64Range(t *testing.T) {
	Check(t, func(t *T, min float64, max float64) {
		Assume(min <= max)
		f := ufloat64FromParts(genUfloatRange(t.src.s, min, max, float64SignifBits))
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64sMin(0), Float64sMin(0))
}

func TestGenFloat32Range(t *testing.T) {
	Check(t, func(t *T, min float32, max float32) {
		Assume(min <= max)
		f := float32FromParts(genFloatRange(t.src.s, float64(min), float64(max), float32SignifBits))
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float32bits(f), min, max, math.Float32bits(min), math.Float32bits(max))
		}
	}, Float32s(), Float32s())
}

func TestGenFloat64Range(t *testing.T) {
	Check(t, func(t *T, min float64, max float64) {
		Assume(min <= max)
		f := float64FromParts(genFloatRange(t.src.s, min, max, float64SignifBits))
		if f < min || f > max {
			t.Fatalf("%v (0x%x) outside of [%v, %v] ([0x%x, 0x%x])", f, math.Float64bits(f), min, max, math.Float64bits(min), math.Float64bits(max))
		}
	}, Float64s(), Float64s())
}
