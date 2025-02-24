// Copyright 2022 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"context"
	"testing"

	. "pgregory.net/rapid"
)

func checkInt(t *T) {
	answer := Int().Draw(t, "answer")
	if answer == 42 {
		t.Fatalf("fuzzing works")
	}
}

func checkSlice(t *T) {
	slice := SliceOfN(Int(), 5, 5).Draw(t, "slice")
	if slice[0] < slice[1] && slice[1] < slice[2] && slice[2] < slice[3] && slice[3] < slice[4] {
		t.Fatalf("fuzzing works")
	}
}

func checkString(t *T) {
	hello := String().Draw(t, "hello")
	if hello == "world" {
		t.Fatalf("fuzzing works")
	}
}

func checkStuckStateMachine(t *T) {
	die := 0
	t.Repeat(map[string]func(*T){
		"roll": func(t *T) {
			if die == 6 {
				t.Skip("game over")
			}
			die = IntRange(1, 6).Draw(t, "die")
		},
	})
}

func TestRapidInt(t *testing.T) {
	t.Skip()
	Check(t, checkInt)
}
func TestRapidSlice(t *testing.T) {
	t.Skip()
	Check(t, checkSlice)
}
func TestRapidString(t *testing.T) {
	t.Skip()
	Check(t, checkString)
}
func TestRapidStuckStateMachine(t *testing.T) {
	t.Skip()
	Check(t, checkStuckStateMachine)
}

func FuzzInt(f *testing.F)               { f.Fuzz(MakeFuzz(checkInt)) }
func FuzzSlice(f *testing.F)             { f.Fuzz(MakeFuzz(checkSlice)) }
func FuzzString(f *testing.F)            { f.Fuzz(MakeFuzz(checkString)) }
func FuzzStuckStateMachine(f *testing.F) { f.Fuzz(MakeFuzz(checkStuckStateMachine)) }

func FuzzContext(f *testing.F) {
	type key struct{}

	var ctx context.Context
	f.Fuzz(MakeFuzz(func(t *T) {
		// Assign to outer variable
		// so we can check it after the fuzzing.
		ctx = context.WithValue(t.Context(), key{}, "value")
		if err := ctx.Err(); err != nil {
			t.Fatalf("context must be valid: %v", err)
		}
	}))

	// ctx is set only if the fuzzing function was called.
	if ctx != nil {
		if err := ctx.Err(); err == nil {
			f.Fatalf("context must be canceled")
		}

		if want, got := "value", ctx.Value(key{}); want != got {
			f.Fatalf("context must have value %q, got %q", want, got)
		}
	}
}

func FuzzCleanup(f *testing.F) {
	var state []bool
	f.Fuzz(MakeFuzz(func(t *T) {
		idx := len(state)
		state = append(state, false)
		t.Cleanup(func() {
			state[idx] = true
		})
	}))

	for _, ok := range state {
		if !ok {
			f.Fatalf("cleanup must be called")
		}
	}
}
