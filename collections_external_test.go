// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"reflect"
	"strconv"
	"testing"

	. "github.com/flyingmutant/rapid"
)

func TestSliceOf(t *testing.T) {
	t.Parallel()

	gens := []*Generator{
		SliceOf(Boolean()),
		SliceOf(Byte()),
		SliceOf(Int()),
		SliceOf(Uint()),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) {
			v := g.Draw(t, "v")
			if rv(v).Kind() != reflect.Slice {
				t.Fatalf("got not a slice")
			}
			if rv(v).Len() == 0 {
				t.Skip("empty")
			}
		}))
	}
}

func TestSliceOfDistinct(t *testing.T) {
	t.Parallel()

	g := SliceOfDistinct(Int(), nil)

	Check(t, func(t *T) {
		s := g.Draw(t, "s").([]int)
		m := map[int]struct{}{}
		for _, i := range s {
			m[i] = struct{}{}
		}
		if len(m) != len(s) {
			t.Fatalf("%v unique out of %v", len(m), len(s))
		}
	})
}

func TestSliceOfDistinctBy(t *testing.T) {
	t.Parallel()

	g := SliceOfDistinct(Int(), func(i int) string { return strconv.Itoa(i % 5) })

	Check(t, func(t *T) {
		s := g.Draw(t, "s").([]int)
		m := map[int]struct{}{}
		for _, i := range s {
			m[i%5] = struct{}{}
		}
		if len(m) != len(s) {
			t.Fatalf("%v unique out of %v", len(m), len(s))
		}
	})
}

func TestMapOf(t *testing.T) {
	t.Parallel()

	gens := []*Generator{
		MapOf(Boolean(), Int()),
		MapOf(Int(), Uint()),
		MapOf(Uint(), SliceOf(Boolean())),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) {
			v := g.Draw(t, "v")
			if rv(v).Kind() != reflect.Map {
				t.Fatalf("got not a map")
			}
			if rv(v).Len() == 0 {
				t.Skip("empty")
			}
		}))
	}
}

func TestMapOfValues(t *testing.T) {
	t.Parallel()

	g := MapOfValues(Custom(genStruct), func(s testStruct) int { return s.x })

	Check(t, func(t *T) {
		m := g.Draw(t, "m").(map[int]testStruct)
		for k, v := range m {
			if k != v.x {
				t.Fatalf("got key %v with value %v", k, v)
			}
		}
	})
}

func TestArrayOf(t *testing.T) {
	t.Parallel()

	elems := []*Generator{Boolean(), Int(), Uint()}
	counts := []int{0, 1, 3, 17}

	for _, e := range elems {
		for _, c := range counts {
			g := ArrayOf(c, e)
			t.Run(g.String(), MakeCheck(func(t *T) {
				v := g.Draw(t, "v")
				if rv(v).Len() != c {
					t.Fatalf("len is %v instead of %v", rv(v).Len(), c)
				}
			}))
		}
	}
}

func TestCollectionLenLimits(t *testing.T) {
	t.Parallel()

	genFuncs := []func(i, j int) *Generator{
		func(i, j int) *Generator { return StringOfN(Byte(), i, j, -1) },
		func(i, j int) *Generator { return SliceOfN(Byte(), i, j) },
		func(i, j int) *Generator { return SliceOfNDistinct(Byte(), i, j, nil) },
		func(i, j int) *Generator { return SliceOfNDistinct(Int(), i, j, func(n int) int { return n % j }) },
		func(i, j int) *Generator { return MapOfN(Int(), Int(), i, j) },
		func(i, j int) *Generator { return MapOfNValues(Int(), i, j, nil) },
		func(i, j int) *Generator { return MapOfNValues(Int(), i, j, func(n int) int { return n % j }) },
	}

	for i, gf := range genFuncs {
		t.Run(strconv.Itoa(i), MakeCheck(func(t *T) {
			minLen := IntRange(0, 256).Draw(t, "minLen").(int)
			maxLen := IntMin(minLen).Draw(t, "maxLen").(int)

			s := rv(gf(minLen, maxLen).Draw(t, "s"))
			if s.Len() < minLen {
				t.Fatalf("got collection of length %v with minLen %v", s.Len(), minLen)
			}
			if s.Len() > maxLen {
				t.Fatalf("got collection of length %v with maxLen %v", s.Len(), maxLen)
			}
		}))
	}
}
