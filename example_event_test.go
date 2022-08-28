// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"
	"testing"

	"pgregory.net/rapid"
)

func ExampleEvent(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// For any integers x, y ...
		x := rapid.Int().Draw(t, "x").(int)
		y := rapid.Int().Draw(t, "y").(int)
		// ... report them ...
		rapid.Event(t, "x", fmt.Sprintf("%d", x))
		rapid.Event(t, "y", fmt.Sprintf("%d", y))

		// ... the property holds
		if x+y != y+x {
			t.Fatalf("associativty of + does not hold")
		}
		// statistics are printed after the property (if called with go test -v)
	})
	// Output:
}
