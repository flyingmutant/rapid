// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"flag"
	"math/bits"
	"reflect"
	"sort"
	"strconv"
	"testing"

	. "github.com/flyingmutant/rapid"
)

var (
	flaky = flag.Bool("flaky.ext", false, "run flaky external tests")
	rv    = reflect.ValueOf
)

func TestIntsExamples(t *testing.T) {
	gens := []*Generator{
		Int(),
		IntMin(-3),
		IntMax(3),
		IntRange(-3, 7),
		IntRange(-1000, 1000000),
		IntRange(0, 9),
		IntRange(0, 15),
		IntRange(10, 100),
		IntRange(100, 10000),
		IntRange(100, 1000000),
		IntRange(100, 1<<60-1),
	}

	for _, g := range gens {
		t.Run(g.String(), func(t *testing.T) {
			var vals []int
			for i := 0; i < 100; i++ {
				n, _, _ := g.Example()
				vals = append(vals, int(rv(n).Int()))
			}
			sort.Ints(vals)

			for _, i := range vals {
				t.Log(i)
			}
		})
	}
}

func createGen(ctor interface{}, args ...interface{}) *Generator {
	refArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		refArgs[i] = rv(arg)
	}

	return rv(ctor).Call(refArgs)[0].Interface().(*Generator)
}

func TestIntsMinMaxRange(t *testing.T) {
	t.Parallel()

	data := []struct {
		g      *Generator
		min    interface{}
		max    interface{}
		range_ interface{}
	}{
		{Int(), IntMin, IntMax, IntRange},
		{Int8(), Int8Min, Int8Max, Int8Range},
		{Int16(), Int16Min, Int16Max, Int16Range},
		{Int32(), Int32Min, Int32Max, Int32Range},
		{Int64(), Int64Min, Int64Max, Int64Range},
	}

	for _, d := range data {
		t.Run(d.g.String(), MakeCheck(func(t *T) {
			min := d.g.Draw(t, "min")
			max := d.g.Draw(t, "max")
			if rv(min).Int() > rv(max).Int() {
				t.Skip("min > max")
			}

			i := createGen(d.min, min).Draw(t, "i")
			if rv(i).Int() < rv(min).Int() {
				t.Fatalf("got %v which is less than min %v", i, min)
			}

			j := createGen(d.max, max).Draw(t, "j")
			if rv(j).Int() > rv(max).Int() {
				t.Fatalf("got %v which is more than max %v", j, max)
			}

			k := createGen(d.range_, min, max).Draw(t, "k")
			if rv(k).Int() < rv(min).Int() || rv(k).Int() > rv(max).Int() {
				t.Fatalf("got %v which is out of bounds [%v, %v]", k, min, max)
			}
		}))
	}
}

func TestUintsMinMaxRange(t *testing.T) {
	t.Parallel()

	data := []struct {
		g      *Generator
		min    interface{}
		max    interface{}
		range_ interface{}
	}{
		{Byte(), ByteMin, ByteMax, ByteRange},
		{Uint(), UintMin, UintMax, UintRange},
		{Uint8(), Uint8Min, Uint8Max, Uint8Range},
		{Uint16(), Uint16Min, Uint16Max, Uint16Range},
		{Uint32(), Uint32Min, Uint32Max, Uint32Range},
		{Uint64(), Uint64Min, Uint64Max, Uint64Range},
		{Uintptr(), UintptrMin, UintptrMax, UintptrRange},
	}

	for _, d := range data {
		t.Run(d.g.String(), MakeCheck(func(t *T) {
			min := d.g.Draw(t, "min")
			max := d.g.Draw(t, "max")
			if rv(min).Uint() > rv(max).Uint() {
				t.Skip("min > max")
			}

			i := createGen(d.min, min).Draw(t, "i")
			if rv(i).Uint() < rv(min).Uint() {
				t.Fatalf("got %v which is less than min %v", i, min)
			}

			j := createGen(d.max, max).Draw(t, "j")
			if rv(j).Uint() > rv(max).Uint() {
				t.Fatalf("got %v which is more than max %v", j, max)
			}

			k := createGen(d.range_, min, max).Draw(t, "k")
			if rv(k).Uint() < rv(min).Uint() || rv(k).Uint() > rv(max).Uint() {
				t.Fatalf("got %v which is out of bounds [%v, %v]", k, min, max)
			}
		}))
	}
}

func TestIntsBoundCoverage(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T) {
		min := Int().Draw(t, "min").(int)
		max := Int().Draw(t, "max").(int)
		if min > max {
			min, max = max, min
		}

		g := IntRange(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 250; i++ {
			n_, _, _ := g.Example(uint64(i))
			n := n_.(int)

			if n == min {
				gotMin = true
			}
			if n == max {
				gotMax = true
			}
			if n == 0 {
				gotZero = true
			}
			if gotMin && gotMax && (min > 0 || max < 0 || gotZero) {
				return
			}
		}

		t.Fatalf("[%v, %v]: got min %v, got max %v, got zero %v", min, max, gotMin, gotMax, gotZero)
	})
}

func TestBytesCoverage(t *testing.T) {
	t.Parallel()
	if !*flaky {
		t.Skip("flaky")
	}

	for b := 0; b < 256; b++ {
		t.Run(strconv.Itoa(b), func(t *testing.T) {
			g := Byte().Filter(func(v byte) bool { return v == byte(b) })
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}

func TestIntsCoverage(t *testing.T) {
	t.Parallel()
	if !*flaky {
		t.Skip("flaky")
	}

	filters := []func(int) bool{
		func(i int) bool { return i == 0 },
		func(i int) bool { return i == 1 },
		func(i int) bool { return i == -1 },
		func(i int) bool { return i%2 == 0 },
		func(i int) bool { return i%17 == 0 },
		func(i int) bool { return i > 0 && i < 100 },
		func(i int) bool { return i < 0 && i > -100 },
		func(i int) bool { return i > 1<<30 },
		func(i int) bool { return i < -(1 << 30) },
	}

	if bits.UintSize == 64 {
		filters = append(filters, func(i int) bool { return i > 1<<62 })
		filters = append(filters, func(i int) bool { return i < -(1 << 62) })
	}

	for i, fn := range filters {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := Int().Filter(fn)
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}

func TestUintsCoverage(t *testing.T) {
	t.Parallel()
	if !*flaky {
		t.Skip("flaky")
	}

	filters := []func(uint) bool{
		func(i uint) bool { return i == 0 },
		func(i uint) bool { return i == 1 },
		func(i uint) bool { return i%2 == 0 },
		func(i uint) bool { return i%17 == 0 },
		func(i uint) bool { return i > 0 && i < 100 },
		func(i uint) bool { return i > 1<<31 },
	}

	if bits.UintSize == 64 {
		filters = append(filters, func(i uint) bool { return i > 1<<63 })
	}

	for i, fn := range filters {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := Uint().Filter(fn)
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}
