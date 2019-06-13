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
		Ints(),
		IntsMin(-3),
		IntsMax(3),
		IntsRange(-3, 7),
		IntsRange(-1000, 1000000),
		IntsRange(0, 9),
		IntsRange(0, 15),
		IntsRange(10, 100),
		IntsRange(100, 10000),
		IntsRange(100, 1000000),
		IntsRange(100, 1<<60-1),
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
	data := []struct {
		g      *Generator
		min    interface{}
		max    interface{}
		range_ interface{}
	}{
		{Ints(), IntsMin, IntsMax, IntsRange},
		{Int8s(), Int8sMin, Int8sMax, Int8sRange},
		{Int16s(), Int16sMin, Int16sMax, Int16sRange},
		{Int32s(), Int32sMin, Int32sMax, Int32sRange},
		{Int64s(), Int64sMin, Int64sMax, Int64sRange},
	}

	for _, d := range data {
		t.Run(d.g.String(), MakeCheck(func(t *T) {
			min := d.g.Draw(t, "min")
			max := d.g.Draw(t, "max")
			Assume(rv(min).Int() <= rv(max).Int())

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
	data := []struct {
		g      *Generator
		min    interface{}
		max    interface{}
		range_ interface{}
	}{
		{Bytes(), BytesMin, BytesMax, BytesRange},
		{Uints(), UintsMin, UintsMax, UintsRange},
		{Uint8s(), Uint8sMin, Uint8sMax, Uint8sRange},
		{Uint16s(), Uint16sMin, Uint16sMax, Uint16sRange},
		{Uint32s(), Uint32sMin, Uint32sMax, Uint32sRange},
		{Uint64s(), Uint64sMin, Uint64sMax, Uint64sRange},
		{Uintptrs(), UintptrsMin, UintptrsMax, UintptrsRange},
	}

	for _, d := range data {
		t.Run(d.g.String(), MakeCheck(func(t *T) {
			min := d.g.Draw(t, "min")
			max := d.g.Draw(t, "max")
			Assume(rv(min).Uint() <= rv(max).Uint())

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
	Check(t, func(t *T, min int, max int) {
		Assume(min <= max)

		g := IntsRange(min, max)
		var gotMin, gotMax, gotZero bool
		for i := 0; i < 250; i++ {
			n_, _, _ := g.Example(uint64(i) + 1)
			n := int(rv(n_).Int())

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
	}, Ints(), Ints())
}

func TestBytesCoverage(t *testing.T) {
	if !*flaky {
		t.Skip()
	}

	for b := 0; b < 256; b++ {
		t.Run(strconv.Itoa(b), func(t *testing.T) {
			g := Bytes().Filter(func(v byte) bool { return v == byte(b) })
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}

func TestIntsCoverage(t *testing.T) {
	if !*flaky {
		t.Skip()
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
			g := Ints().Filter(fn)
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}

func TestUintsCoverage(t *testing.T) {
	if !*flaky {
		t.Skip()
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
			g := Uints().Filter(fn)
			_, n, err := g.Example()
			if err != nil {
				t.Errorf("failed to find an example in %v tries: %v", n, err)
			}
		})
	}
}
