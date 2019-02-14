// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"math/bits"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

const shrinkTestRuns = 10

func TestShrink_Trivial(t *testing.T) {
	checkShrink(t, Bind(func(t *T, i int) {
		if i > 1000000 {
			t.Fail()
		}
	}, Ints()), pack(1000001))
}

func TestShrink_NegativeLt(t *testing.T) {
	checkShrink(t, Bind(func(t *T, i int) {
		if i < -1000000 {
			t.Fail()
		}
	}, Ints()), pack(-1000001))
}

func TestShrink_NegativeLe(t *testing.T) {
	checkShrink(t, Bind(func(t *T, i int) {
		if i <= 0 {
			t.Fail()
		}
	}, Ints()), pack(0))
}

func TestShrink_NegativeGt(t *testing.T) {
	checkShrink(t, Bind(func(t *T, i int) {
		if i > -1000000 {
			t.Fail()
		}
	}, Ints()), pack(0))
}

func TestShrink_CollectionElements(t *testing.T) {
	checkShrink(t, Bind(func(t *T, s []int) {
		n := 0
		for _, i := range s {
			if i > 1000000 {
				n++
			}
		}
		if n > 1 {
			t.Fail()
		}
	}, SlicesOf(Ints())), pack([]int{1000001, 1000001}))
}

func TestShrink_CollectionIndex(t *testing.T) {
	checkShrink(t, Bind(func(t *T, s []int) {
		ix := t.Draw(IntsRange(0, len(s)-1), "ix").(int)

		if s[ix] >= 100 {
			t.Fail()
		}

	}, SlicesOfN(Ints(), 1, -1)), pack([]int{100}), 0)
}

func TestShrink_CollectionSpan(t *testing.T) {
	checkShrink(t, Bind(func(t *T, s []int) {
		if len(s)%3 == 1 && s[len(s)-1] >= 100 {
			t.Fail()
		}
	}, SlicesOfN(Ints(), 4, -1)), pack([]int{0, 0, 0, 100}))
}

func TestShrink_Sort(t *testing.T) {
	checkShrink(t, Bind(func(t *T, s []int) {
		sort.Ints(s)
		last := 0
		for _, i := range s {
			if i == last {
				return
			}
			last = i
		}
		t.Fail()
	}, SlicesOfN(IntsMin(1), 5, -1)), pack([]int{1, 2, 3, 4, 5}))
}

func TestShrink_Strings(t *testing.T) {
	checkShrink(t, Bind(func(t *T, s1 string, s2 string) {
		if len(s1) > len(s2) {
			t.Fail()
		}
	}, Strings(), Strings()), pack("?", ""))
}

func TestMinimize_UnsetBits(t *testing.T) {
	Check(t, func(t *T, mask uint64) {
		best := minimize(math.MaxUint64, func(x uint64, s string) bool { return x&mask == mask })
		if best != mask {
			t.Fatalf("unset to %v instead of %v", bin(best), bin(mask))
		}
	}, Uint64sRange(0, math.MaxUint64))
}

func TestMinimize_SortBits(t *testing.T) {
	Check(t, func(t *T, u uint64) {
		n := bits.OnesCount64(u)
		v := uint64(1<<uint(n) - 1)

		best := minimize(u, func(x uint64, s string) bool { return bits.OnesCount64(x) == n })
		if best != v {
			t.Fatalf("minimized to %v instead of %v (%v bits set)", bin(best), bin(v), n)
		}
	}, Uint64sRange(0, math.MaxUint64))
}

func TestMinimize_LowerBound(t *testing.T) {
	Check(t, func(t *T) {
		min := t.Draw(Uint64s(), "min").(uint64)
		u := t.Draw(Uint64sMin(min), "u").(uint64)

		best := minimize(u, func(x uint64, s string) bool { return x >= min })
		if best != min {
			t.Fatalf("found %v instead of %v", best, min)
		}
	})
}

func checkShrink(t *testing.T, prop func(*T), draws ...Value) {
	t.Helper()

	for i := 0; i < shrinkTestRuns; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, _, buf, err1, err2 := doCheck(t, prop)
			if err1 == nil && err2 == nil {
				t.Fatalf("shrink test did not fail")
			}
			if traceback(err1) != traceback(err2) {
				t.Fatalf("flaky shrink test\nTraceback (%v):\n%v\nOriginal traceback (%v):\n%v", err2, traceback(err2), err1, traceback(err1))
			}

			_ = checkOnce(newT(t, newBufBitStream(buf, false), false, draws...), prop)
		})
	}
}

func pack(fields ...Value) Value {
	vals := make([]reflect.Value, len(fields))
	typs := make([]reflect.Type, len(fields))

	for i, field := range fields {
		vals[i] = reflect.ValueOf(field)
		typs[i] = vals[i].Type()
	}

	return packTuple(tupleOf(typs), vals...).Interface()
}

func bin(u uint64) string {
	return "0b" + strconv.FormatUint(u, 2)
}
