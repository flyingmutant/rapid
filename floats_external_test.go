// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"testing"

	. "github.com/flyingmutant/rapid"
)

func TestFloat32sExamples(t *testing.T) {
	g := Float32s()

	for i := 0; i < 100; i++ {
		f, _, _ := g.Example()
		t.Log(f)
	}
}

func TestFloat64sExamples(t *testing.T) {
	g := Float64s()

	for i := 0; i < 100; i++ {
		f, _, _ := g.Example()
		t.Log(f)
	}
}
