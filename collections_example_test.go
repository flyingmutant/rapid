// Copyright 2020 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"
	"math"

	"pgregory.net/rapid"
)

func ExampleSliceOf() {
	gen := rapid.SliceOf(rapid.IntRange(math.MinInt32, math.MaxInt32))

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// [1 -366 7 -236 14 -49 -65303 -1 25522 -2 134 -951504605 -2147483648 17690368 -2 8814 10 -4960809 0 -11]
	// [-3 -2 -1 -3 -25381941 -1 -2 -661644976]
	// [0 308 -2 5 -211 3 1 14 22415068 -11]
	// [78 -3 -15 -1 0 -1 332 6 -1 -1440]
	// []
}

func ExampleSliceOfN() {
	gen := rapid.SliceOfN(rapid.IntRange(math.MinInt32, math.MaxInt32), 5, 5)

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// [1 -366 7 -236 14]
	// [-3 -2 -1 -3 -25381941]
	// [0 308 -2 5 -211]
	// [78 -3 -15 -1 0]
	// [15402512 14 2131 -631093255 -3]
}

func ExampleSliceOfDistinct() {
	gen := rapid.SliceOfDistinct(rapid.IntRange(0, math.MaxInt32), func(i int) int { return i % 2 })

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// [1]
	// [2 1]
	// [0 1]
	// [78]
	// []
}

func ExampleSliceOfNDistinct() {
	gen := rapid.SliceOfNDistinct(rapid.IntRange(0, math.MaxInt32), 2, 2, func(i int) int { return i % 2 })

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// [25522 1]
	// [2 1]
	// [0 1]
	// [0 1439]
	// [15402512 2131]
}

func ExampleMapOf() {
	gen := rapid.MapOf(rapid.IntRange(math.MinInt32, math.MaxInt32), rapid.StringMatching(`[a-z]+`))

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// map[1:nhlgqwasbggbaociac 4804:r]
	// map[-168:pizpv -3:bacuabp 0:bi]
	// map[-235083557:gewf -2132:b -33:a -4:b -2:v 0:ubfsdbowrja 82:braigey 142:lvcprss 239:tcozav]
	// map[-30:h 78:coaaamcasnapgaad]
	// map[]
}

func ExampleMapOfN() {
	gen := rapid.MapOfN(rapid.IntRange(math.MinInt32, math.MaxInt32), rapid.StringMatching(`[a-z]+`), 5, 5)

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// map[-124912695:bd -423:bbdbcs 1:nhlgqwasbggbaociac 15:kmdnpmcbuagzr 4804:r]
	// map[-235476:d -168:pizpv -3:bacuabp 0:bi 7713:rzkneb]
	// map[-235083557:gewf -2:v 0:ubfsdbowrja 142:lvcprss 239:tcozav]
	// map[-2488147:j -9267:aafmd -236:o -30:h 78:coaaamcasnapgaad]
	// map[-378:y 0:paai 2:b 15402512:otg 24810092:qign]
}

func ExampleMapOfValues() {
	gen := rapid.MapOfValues(rapid.StringMatching(`[a-z]+`), func(s string) int { return len(s) })

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// map[2:dr 7:xguehfc 11:sbggbaociac]
	// map[2:bp 5:jarxz 6:ebzkwa]
	// map[1:j 2:aj 3:gjl 4:vayt 5:eeeqa 6:riacaa 7:stcozav 8:mfdhbzcz 9:fxmcadagf 10:bgsbraigey 15:gxongygnxqlovib]
	// map[2:ub 8:waraafmd 10:bfiqcaxazu 16:rjgqimcasnapgaad 17:gckfbljafcedhcvfc]
	// map[]
}

func ExampleMapOfNValues() {
	gen := rapid.MapOfNValues(rapid.StringMatching(`[a-z]+`), 5, 5, func(s string) int { return len(s) })

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// map[1:s 2:dr 3:anc 7:xguehfc 11:sbggbaociac]
	// map[1:b 2:bp 4:ydag 5:jarxz 6:ebzkwa]
	// map[1:j 3:gjl 5:eeeqa 7:stcozav 9:fxmcadagf]
	// map[2:ub 8:waraafmd 10:bfiqcaxazu 16:rjgqimcasnapgaad 17:gckfbljafcedhcvfc]
	// map[1:k 2:ay 3:wzb 4:dign 7:faabhcb]
}
