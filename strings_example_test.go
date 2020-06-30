// Copyright 2020 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// String generation depends on the Unicode tables, which change with Go versions:
// +build go1.14

package rapid_test

import (
	"fmt"
	"unicode"

	"pgregory.net/rapid"
)

func ExampleRune() {
	gen := rapid.Rune()

	for i := uint64(0); i < 25; i++ {
		if i%5 == 0 {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
		fmt.Printf("%q", gen.Example(i))
	}
	// Output:
	// '\\' '\ufeff' '?' '~' '-'
	// '0' '$' '!' '`' '\ue05d'
	// '"' '&' '#' '\u0604' 'A'
	// '&' '茞' '@' '#' '|'
	// '⊙' '𝩔' '$' '҈' '\r'
}

func ExampleRuneFrom() {
	gens := []*rapid.Generator{
		rapid.RuneFrom([]rune{'A', 'B', 'C'}),
		rapid.RuneFrom(nil, unicode.Cyrillic, unicode.Greek),
		rapid.RuneFrom([]rune{'⌘'}, &unicode.RangeTable{
			R32: []unicode.Range32{{0x1F600, 0x1F64F, 1}},
		}),
	}

	for _, gen := range gens {
		for i := uint64(0); i < 5; i++ {
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

	for i := uint64(0); i < 5; i++ {
		fmt.Printf("%q\n", gen.Example(i))
	}
	// Output:
	// "\\߾⃝!/?Ⱥ֍"
	// "\u2006𑨷"
	// "?﹩\u0603ᾢ"
	// ".*%:<%৲"
	// ""
}
