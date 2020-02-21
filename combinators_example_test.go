// Copyright 2020 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"

	"pgregory.net/rapid"
)

func ExampleCustom() {
	type point struct {
		x int
		y int
	}

	gen := rapid.Custom(func(t *rapid.T) point {
		return point{
			x: rapid.Int().Draw(t, "x").(int),
			y: rapid.Int().Draw(t, "y").(int),
		}
	})

	for i := uint64(0); i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// {-3 1303}
	// {-186981 -59881619}
	// {4 441488606}
	// {-2 -5863986}
	// {43 -3513}
}

func ExampleJust() {
	gen := rapid.Just(42)

	for i := uint64(0); i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// 42
	// 42
	// 42
	// 42
	// 42
}

func ExampleSampledFrom() {
	gen := rapid.SampledFrom([]int{1, 2, 3})

	for i := uint64(0); i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// 2
	// 3
	// 2
	// 3
	// 1
}

func ExampleOneOf() {
	gen := rapid.OneOf(rapid.IntRange(1, 10), rapid.IntRange(100, 1000))

	for i := uint64(0); i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// 159
	// 10
	// 109
	// 2
	// 9
}

func ExamplePtr() {
	gen := rapid.Ptr(rapid.Int(), true)

	for i := uint64(0); i < 5; i++ {
		v := gen.Example(i).(*int)
		if v == nil {
			fmt.Println("<nil>")
		} else {
			fmt.Println("(*int)", *v)
		}
	}
	// Output:
	// (*int) 1
	// (*int) -3
	// <nil>
	// (*int) 590
	// <nil>
}
