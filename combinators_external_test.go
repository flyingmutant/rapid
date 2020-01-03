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

func genBool(t *T) bool {
	return Boolean().Draw(t, "").(bool)
}

func genSlice(t *T) []uint64 {
	return []uint64{
		Uint64().Draw(t, "").(uint64),
		Uint64().Draw(t, "").(uint64),
	}
}

func genStruct(t *T) testStruct {
	return testStruct{
		x: Int().Draw(t, "x").(int),
		y: Int().Draw(t, "y").(int),
	}
}

func TestCustom(t *testing.T) {
	t.Parallel()

	gens := []*Generator{
		Custom(genBool),
		Custom(genSlice),
		Custom(genStruct),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) { g.Draw(t, "") }))
	}
}

func TestFilter(t *testing.T) {
	t.Parallel()

	g := Int().Filter(func(i int) bool { return i >= 0 })

	Check(t, func(t *T) {
		v := g.Draw(t, "v").(int)
		if v < 0 {
			t.Fatalf("got negative %v", v)
		}
	})
}

func TestMap(t *testing.T) {
	t.Parallel()

	g := Int().Map(strconv.Itoa)

	Check(t, func(t *T) {
		s := g.Draw(t, "s").(string)
		_, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("Atoi() error %v", err)
		}
	})
}

func TestSampledFrom(t *testing.T) {
	t.Parallel()

	gens := []*Generator{
		Just(3),
		SampledFrom([]int{3, 5, 7}),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) {
			n := g.Draw(t, "n").(int)
			if n != 3 && n != 5 && n != 7 {
				t.Fatalf("got impossible %v", n)
			}
		}))
	}
}

func TestOneOf(t *testing.T) {
	t.Parallel()

	pos := Int().Filter(func(v int) bool { return v >= 10 })
	neg := Int().Filter(func(v int) bool { return v <= -10 })
	g := OneOf(pos, neg)

	Check(t, func(t *T) {
		n := g.Draw(t, "n").(int)
		if n > -10 && n < 10 {
			t.Fatalf("got impossible %v", n)
		}
	})
}

func TestPtr(t *testing.T) {
	t.Parallel()

	for _, allowNil := range []bool{false, true} {
		t.Run(fmt.Sprintf("allowNil=%v", allowNil), MakeCheck(func(t *T) {
			i := Ptr(Int(), allowNil).Draw(t, "i").(*int)
			if i == nil && !allowNil {
				t.Fatalf("got nil pointer")
			}
		}))
	}
}
