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

const floatExampleFormat = "%30g %10.3g % 5d % 20d % 16x"

func TestFloat32sExamples(t *testing.T) {
	g := Float32s()

	var vals []float64
	for i := 0; i < 100; i++ {
		f_, _, _ := g.Example()
		vals = append(vals, float64(f_.(float32)))
	}
	sort.Float64s(vals)

	for _, f := range vals {
		t.Logf(floatExampleFormat, f, f, int(math.Log10(math.Abs(f))), int64(f), math.Float32bits(float32(f)))
	}
}

func TestFloat64sExamples(t *testing.T) {
	g := Float64s()

	var vals []float64
	for i := 0; i < 100; i++ {
		f_, _, _ := g.Example()
		vals = append(vals, f_.(float64))
	}
	sort.Float64s(vals)

	for _, f := range vals {
		t.Logf(floatExampleFormat, f, f, int(math.Log10(math.Abs(f))), int64(f), math.Float64bits(f))
	}
}
