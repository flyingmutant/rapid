// Copyright 2022 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
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

func FuzzInt(f *testing.F)    { f.Fuzz(MakeFuzz(checkInt)) }
func FuzzSlice(f *testing.F)  { f.Fuzz(MakeFuzz(checkSlice)) }
func FuzzString(f *testing.F) { f.Fuzz(MakeFuzz(checkString)) }
