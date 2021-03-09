// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// String generation depends on the Unicode tables, which change with Go versions:
// +build go1.16

package rapid_test

import (
	"fmt"
	"net/url"

	"pgregory.net/rapid"
)

func ExampleDomain() {
	gen := rapid.Domain()

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}

	// Output:
	// D1C.TRaVElErs
	// C.cuISiNeLlA
	// r.abbVIe
	// MC0zJ.aCcOuntAnTs
	// T6hFdv10.aaa
}

func ExampleURL() {
	gen := rapid.URL()

	for i := 0; i < 5; i++ {
		e := gen.Example(i).(url.URL)
		fmt.Println(e.String())
	}

	// Output:
	// https://[e506:816b:407:316:fb4c:ffa0:e208:dc0e]/%F0%97%B0%A0%F0%92%91%9CX/1=%22?%C4%90%F0%90%A9%87%26#%F0%91%B0%B31%CC%88%CC%81D
	// http://G.BLoG/%E0%AD%8C~%F0%9F%AF%8A%22%D6%93%E0%A9%AD%E1%B3%930%D0%8A/%C2%BC%E0%B4%BC3%F0%9D%9F%B9%F0%91%86%B6%C2%B2%E0%B3%A9%CC%80D/%7C+%F0%9F%81%A6+%5D%CC%81%CB%85/%CC%80/%E1%B0%BF/%CD%82K%E0%A5%A9%CC%81#%CC%82
	// https://1.47.4.5:11/+%3E%E2%9F%BCK//A%DB%97%F0%90%AC%BB$%F0%91%97%97%E2%82%A5%F0%9D%A8%8D%E0%BC%95%E0%B5%A8%3C%E0%BE%AE%F0%97%8A%B1%E2%9E%8E%E0%B9%91$v%CC%80/%CC%94Z%E4%87%94?%F0%91%91%8EC%C2%B9%E2%8A%A5%F0%91%91%8D1%E0%A0%BE%CB%BE%C3%9D%E1%B3%A4%E0%AB%A6%CC%81%CC%86&%E2%A4%88%F0%91%88%BC%24B%F0%96%BA%90%CC%9A&&%CC%80%C2%A7%E8%93%8B&#%E0%AB%AE%F0%92%91%91
	// http://G.hM/%CC%80%E0%A0%B1%CC%82%CC%80%F0%9E%A5%9F/:%21J%E2%9D%87#L%CC%82%E9%98%8C%22
	// http://1.1.4.6:2/%F0%A7%A8%A4%F0%A1%AD%8D%E2%92%8B0/%DC%B4B?%E2%80%A60%CC%80%C3%B7&%CC%81%CC%A2%21%E0%AF%AB%CC%81%C3%A4&%F0%9F%AF%8A%EA%99%AF%CC%80%E0%A5%AD&%E5%8B%B71&%E1%B7%8F%CC%8B%E2%87%94%E2%90%8E%F3%A0%87%97%2B%C4%B9%C5%A9%5E&B%F0%A8%83%A2%E1%89%AE%E0%A0%BD%F0%9D%83%90%C6%81%211A3&%E1%81%8F%23#%CC%80%E0%BF%8B+$
}
