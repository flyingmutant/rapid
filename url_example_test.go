// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// String generation depends on the Unicode tables, which change with Go versions:
//go:build go1.16
// +build go1.16

package rapid_test

import (
	"fmt"

	"pgregory.net/rapid"
)

func ExampleDomain() {
	gen := rapid.Domain()

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}

	// Output:
	// l.UOL
	// S.CYMRU
	// e.ABBVIE
	// x.ACCOUNTANTS
	// R.AAA
}

func ExampleURL() {
	gen := rapid.URL()

	for i := 0; i < 5; i++ {
		e := gen.Example(i)
		fmt.Println(e.String())
	}

	// Output:
	// https://[e506:816b:407:316:fb4c:ffa0:e208:dc0e]/%25F0%2597%25B0%25A0%25F0%2592%2591%259CX/1=%2522?%C4%90%F0%90%A9%87%26#%F0%96%AC%B21%CC%88%CC%81D
	// http://Z.BLOOMBERG:2/%25E2%259C%25A1/1%25F0%2591%2588%25BD/%25F0%259F%25AF%258A%2522%25D6%2593%25E0%25A9%25AD%25E1%25B3%25930%25D0%258A/%25C2%25BC%25E0%25B4%25BC3%25F0%259D%259F%25B9%25F0%2591%2591%2582%25C2%25B2%25E0%25B3%25A9%25CC%2580D/%257C+%25F0%259F%2582%2592+%255D%25CC%2581%25CB%2585/%25CC%2580/%25E1%25B0%25BF/%25CD%2582K%25E0%25A5%25A9%25CC%2581#%CC%82
	// https://1.47.4.5:11/+%253E%25E2%259F%25BCK//A%25DB%2597%25F0%2591%2599%2583$%25E0%25A0%25BD%25E2%2582%25A5%25F0%259D%25A9%25A9%25E0%25BC%2595%25E0%25B5%25A8%253C%25E0%25BE%25AE%25F0%2597%258A%25B1%25E2%259E%258E%25E0%25B9%2591$v%25CC%2580/%25CC%2594Z%25E4%2587%2594?%F0%96%A9%AEC%C2%B9%E2%8A%A5%F0%92%91%B41%E0%A0%BE%CB%BE%C3%9D%E1%B3%A4%E0%AB%A6%CC%81%CC%86&%E2%A4%88%F0%91%BF%BF%24B%F0%96%BA%90%CC%9A&&%CC%80%C2%A7%E8%93%8B&#%E0%AB%AE%F0%92%91%91
	// http://J.HOMESENSE/%25F0%259B%2589%259D%25C2%25B86%25CC%2580%25F0%259E%25A5%259F/:%2521J%25E2%259D%2587#L%CC%82%E9%98%8C%22
	// http://1.1.4.6:2/%25F0%25A7%25A8%25A4%25F0%25A1%25AD%258D%25E2%2592%258B0/%25DC%25B4B?%E2%80%A60%CC%80%C3%B7&%CC%81%CC%A2%21%E0%AF%AB%CC%81%C3%A4&%F0%9F%AF%8A%EA%99%AF%CC%80%E0%A5%AD&%E5%8B%B71&%E1%B7%8F%CC%8B%E2%87%94%E2%90%8E%EA%A3%A0%E0%B5%9A%3D%E5%8E%8A%D9%AAB%F0%A8%83%A2%EF%B8%B4%E0%A0%BD%F0%9D%84%86%C6%81%211A3&%E1%81%8F%23#%CC%80%E0%BF%8B+$
}