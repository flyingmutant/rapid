// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/flyingmutant/rapid"
)

type testStruct struct {
	x int
	y int
}

func genBool(src Source) bool {
	return Booleans().Draw(src, "").(bool)
}

func genSlice(src Source) []uint64 {
	return []uint64{
		Uint64s().Draw(src, "").(uint64),
		Uint64s().Draw(src, "").(uint64),
	}
}

func genStruct(src Source) testStruct {
	return testStruct{
		x: Ints().Draw(src, "x").(int),
		y: Ints().Draw(src, "y").(int),
	}
}

func genPair(src Source) (int, int) {
	return Ints().Draw(src, "").(int), Ints().Draw(src, "").(int)
}

func TestCustom(t *testing.T) {
	gens := []*Generator{
		Custom(genBool),
		Custom(genSlice),
		Custom(genStruct),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) { g.Draw(t, "") }))
	}
}

func TestTupleHoldover(t *testing.T) {
	g := Tuple(Bytes(), Ints()).Map(func(b byte, i int) bool { return i > int(b) })

	Check(t, func(*T, bool) {}, g)
}

func TestTupleUnpackArgs(t *testing.T) {
	g := Custom(genPair).
		Filter(func(x int, y int) bool { return x != y }).
		Map(func(x int, y int) (int, int, int) { return x, x * 3, y * 3 })

	Check(t, func(t *T, a int, b int, c int) {
		if b != a*3 || b == c {
			t.Fatalf("got impossible %v, %v, %v", a, b, c)
		}
	}, g)
}

func TestTupleUnpackDraw(t *testing.T) {
	g := Custom(genPair).Map(func(x int, y int) (int, string) { return x, strconv.Itoa(x) })

	Check(t, func(t *T) {
		var a int
		var b string
		g.Draw(t, "pair", &a, &b)
		if strconv.Itoa(a) != b {
			t.Fatalf("got impossible %v, %v", a, b)
		}
	})
}

func TestTupleCompatibility(t *testing.T) {
	g := MapsOfNValues(OneOf(Custom(genPair), Custom(genPair), Custom(genPair)), 10, -1, nil)

	Check(t, func(t *T) { g.Draw(t, "") })
}

func TestFilter(t *testing.T) {
	g := Ints().Filter(func(i int) bool { return i >= 0 })

	Check(t, func(t *T, v int) {
		if v < 0 {
			t.Fatalf("got negative %v", v)
		}
	}, g)
}

func TestMap(t *testing.T) {
	g := Ints().Map(strconv.Itoa)

	Check(t, func(t *T, s string) {
		_, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("Atoi() error %v", err)
		}
	}, g)
}

func TestSampledFrom(t *testing.T) {
	gens := []*Generator{
		Just(3),
		SampledFrom([]int{3, 5, 7}),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T, n int) {
			if n != 3 && n != 5 && n != 7 {
				t.Fatalf("got impossible %v", n)
			}
		}, g))
	}
}

func TestOneOf(t *testing.T) {
	pos := Ints().Filter(func(v int) bool { return v >= 10 })
	neg := Ints().Filter(func(v int) bool { return v <= -10 })
	g := OneOf(pos, neg)

	Check(t, func(t *T, n int) {
		if n > -10 && n < 10 {
			t.Fatalf("got impossible %v", n)
		}
	}, g)
}

func TestPtrs(t *testing.T) {
	for _, allowNil := range []bool{false, true} {
		t.Run(fmt.Sprintf("allowNil=%v", allowNil), MakeCheck(func(t *T, i *int) {
			if i == nil && !allowNil {
				t.Fatalf("got nil pointer")
			}
		}, Ptrs(Ints(), allowNil)))
	}
}
