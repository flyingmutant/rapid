// Copyright 2023 Gregory Petrosyan <pgregory@pgregory.net>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"
	"unicode"

	"pgregory.net/rapid"
)

func ExampleRune() {
	gen := rapid.Rune()

	for i := 0; i < 25; i++ {
		if i%5 == 0 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
		fmt.Printf("%q", gen.Example(i))
	}
	// Output:
	// '\n' '\x1b' 'A' 'a' '*'
	// '0' '@' '?' '\'' '\ue05d'
	// '<' '%' '!' '\u0604' 'A'
	// '%' '╷' '~' '!' '/'
	// '\u00ad' '𝅾' '@' '҈' ' '
}

func ExampleRuneFrom() {
	gens := []*rapid.Generator[rune]{
		rapid.RuneFrom([]rune{'A', 'B', 'C'}),
		rapid.RuneFrom(nil, unicode.Cyrillic, unicode.Greek),
		rapid.RuneFrom([]rune{'⌘'}, &unicode.RangeTable{
			R32: []unicode.Range32{{0x1F600, 0x1F64F, 1}},
		}),
	}

	for _, gen := range gens {
		for i := 0; i < 5; i++ {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%q", gen.Example(i))
		}
		fmt.Println()
	}
	// Output:
	// 'A' 'A' 'A' 'B' 'A'
	// 'Ͱ' 'Ѥ' 'Ͱ' 'ͱ' 'Ϳ'
	// '😀' '⌘' '😀' '😁' '😋'
}

func ExampleString() {
	gen := rapid.String()

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "\n߾⃝?\rA�֍"
	// "\u2006𑨳"
	// "A＄\u0603ᾢ"
	// "+^#.[#৲"
	// ""
}

func ExampleStringOf() {
	gen := rapid.StringOf(rapid.RuneFrom(nil, unicode.Tibetan))

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "༁༭༇ཬ༆༐༖ༀྸ༁༆༎ༀ༁ཱི༂༨ༀ༂"
	// "༂༁ༀ༂༴ༀ༁ྵ"
	// "ༀ༴༁༅ན༃༁༎ྼ༄༽"
	// "༎༂༎ༀༀༀཌྷ༂ༀྥ"
	// ""
}

func ExampleStringN() {
	gen := rapid.StringN(5, 5, -1)

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "\n߾⃝?\r"
	// "\u2006𑨳#`\x1b"
	// "A＄\u0603ᾢÉ"
	// "+^#.["
	// ".A<a¤"
}

func ExampleStringOfN() {
	gen := rapid.StringOfN(rapid.RuneFrom(nil, unicode.ASCII_Hex_Digit), 6, 6, -1)

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "1D7B6a"
	// "2102e0"
	// "0e15c3"
	// "E2E000"
	// "aEd623"
}

func ExampleStringMatching() {
	gen := rapid.StringMatching(`\(?([0-9]{3})\)?([ .-]?)([0-9]{3})([ .-]?)([0-9]{4})`)

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "(532) 649-9610"
	// "901)-5783983"
	// "914.444.1575"
	// "(316 696.3584"
	// "816)0861080"
}

func ExampleSliceOfBytesMatching() {
	gen := rapid.SliceOfBytesMatching(`[CAGT]+`)

	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "CCTTGAGAGCGATACGGAAG"
	// "GCAGAACT"
	// "AACCGTCGAG"
	// "GGGAAAAGAT"
	// "AGTG"
}
