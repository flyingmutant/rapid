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
		Float32sRange(0, 1),
		Float32sRange(1, 2.5),
		Float32sRange(0, 100),
		Float64s(),
		Float64sMin(-0.1),
		Float64sMin(1),
		Float64sMax(0.1),
		Float64sMax(2.5),
		Float64sRange(0, 1),
		Float64sRange(1, 2.5),
		Float64sRange(0, 100),
	}

	for _, g := range gens {
		t.Run(g.String(), func(t *testing.T) {
			var vals []float64
			for i := 0; i < 100; i++ {
				f, _, _ := g.Example()
				vals = append(vals, rv(f).Float())
			}
			sort.Float64s(vals)

			for _, f := range vals {
				t.Logf("%30g %10.3g % 5d % 20d % 16x", f, f, int(math.Log10(math.Abs(f))), int64(f), math.Float64bits(f))
			}
		})
	}
}
