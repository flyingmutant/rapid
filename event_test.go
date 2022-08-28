// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"regexp"
	"testing"
)

func TestEventEmitter(t *testing.T) {
	t.Parallel()
	Check(t, func(te *T) {
		// te.rawLog.SetOutput(te.rawLog.Output())
		Event(te, "var", "x")
		Event(te, "var", "y")

		// checkMatch(te, fmt.Sprintf("Statistics.*%s", te.Name()), te.output[0])
		// checkMatch(te, "of 2 ", te.output[1])
		// checkMatch(te, "x: 1 \\(50.0+ %", te.output[3])
		// checkMatch(te, "y: 1 \\(50.0+ %", te.output[4])

	})
}

func checkMatch(t *T, pattern, str string) {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		t.Fatalf("Regex compile failed")
	}
	if !matched {
		t.Fatalf("Pattern <%s> does not match in <%s>", pattern, str)
	}
}

func TestTrivialPropertyWithEvents(t *testing.T) {
	t.Parallel()
	Check(t, func(te *T) {
		x := Uint8().Draw(te, "x").(uint8)
		Event(te, "x", fmt.Sprintf("%d", x))
		if x > 255 {
			t.Fatalf("x should fit into a byte")
		}
	})
}

func TestTrivialPropertyWithNumEvents(t *testing.T) {
	t.Parallel()
	Check(t, func(te *T) {
		x := Uint8().Draw(te, "x").(uint8)
		Event(te, "x", fmt.Sprintf("%d", x))
		NumericEvent(te, "x", x)
		if x > 255 {
			t.Fatalf("x should fit into a byte")
		}
	})
}

func TestFilteredGenWithAutoEvents(t *testing.T) {
	t.Parallel()
	Check(t, func(te *T) {
		x := Uint8().Filter(func(n uint8) bool { return n%2 == 0 }).Draw(te, "even x").(uint8)
		Event(te, "x", fmt.Sprintf("%d", x))
		if x > 255 {
			t.Fatalf("x should fit into a byte")
		}
	})
}
