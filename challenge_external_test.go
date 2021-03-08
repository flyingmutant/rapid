// Copyright 2021 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Implementation of shrinking challenges from https://github.com/jlink/shrinking-challenge

package rapid_test

import (
	"reflect"
	"testing"

	. "pgregory.net/rapid"
)

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/bound5.md
func TestChallengeBound5(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOfN(SliceOf(Int16()).Filter(func(s []int16) bool {
			return int16sum(s) < 256
		}), 5, 5).Draw(t, "x").([][]int16)

		n := int16(0)
		for _, e := range x {
			n += int16sum(e)
		}
		if n >= 5*256 {
			t.Fatal()
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/large_union_list.md
func TestChallengeLargeUnionList(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOf(SliceOf(Int())).Draw(t, "x").([][]int)

		m := map[int]bool{}
		for _, s := range x {
			for _, e := range s {
				m[e] = true
				if len(m) > 4 {
					t.Fatal()
				}
			}
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/reverse.md
func TestChallengeReverse(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOf(Int()).Draw(t, "x").([]int)

		if !reflect.DeepEqual(x, reversed(x)) {
			t.Fatal()
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/lengthlist.md
func TestChallengeLengthList(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		n := IntRange(1, 100).Draw(t, "n").(int)
		s := make([]int, n)
		for i := 0; i < n; i++ {
			s[i] = Int().Draw(t, "e").(int)
		}

		for _, e := range s {
			if e >= 900 {
				t.Fatal()
			}
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/difference.md
func TestChallengeDifference1(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := IntMin(1).Draw(t, "x").(int)
		y := IntMin(1).Draw(t, "y").(int)

		if x >= 10 && x == y {
			t.Fatal()
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/difference.md
func TestChallengeDifference2(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := IntMin(1).Draw(t, "x").(int)
		y := IntMin(1).Draw(t, "y").(int)

		if x >= 10 && abs(x-y) >= 1 && abs(x-y) <= 4 {
			t.Fatal()
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/difference.md
func TestChallengeDifference3(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := IntMin(1).Draw(t, "x").(int)
		y := IntMin(1).Draw(t, "y").(int)

		if x >= 10 && abs(x-y) == 1 {
			t.Fatal()
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/coupling.md
func TestChallengeCoupling(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOf(IntRange(0, 10)).Filter(func(s []int) bool {
			for _, n := range s {
				if n >= len(s) {
					return false
				}
			}
			return true
		}).Draw(t, "x").([]int)

		for i, j := range x {
			if i != j && x[j] == i {
				t.Fatal()
			}
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/deletion.md
func TestChallengeDeletion(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOfN(Int(), 1, -1).Draw(t, "x").([]int)
		n := IntRange(0, len(x)-1).Draw(t, "n").(int)

		r := x[n]
		y := append(x[:n], x[n+1:]...)
		for _, e := range y {
			if e == r {
				t.Fatal()
			}
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/distinct.md
func TestChallengeDistinct(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOf(Int()).Draw(t, "x").([]int)

		m := map[int]bool{}
		for _, e := range x {
			m[e] = true
			if len(m) >= 3 {
				t.Fatal()
			}
		}
	})
}

// https://github.com/jlink/shrinking-challenge/blob/main/challenges/nestedlists.md
func TestChallengeNestedLists(t *testing.T) {
	t.Skip("shrinking challenge test is expected to fail")

	Check(t, func(t *T) {
		x := SliceOf(SliceOf(IntRange(0, 0))).Draw(t, "x").([][]int)

		n := 0
		for _, s := range x {
			n += len(s)
		}
		if n > 10 {
			t.Fatal()
		}
	})
}

func int16sum(s []int16) int16 {
	n := int16(0)
	for _, i := range s {
		n += i
	}
	return n
}

func reversed(s []int) []int {
	r := make([]int, len(s))
	for i, n := range s {
		r[len(r)-1-i] = n
	}
	return r
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
