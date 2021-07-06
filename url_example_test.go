// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
	// https://[e506:816b:407:316:fb4c:ffa0:e208:dc0e]/%F0%97%B3%80%F0%92%91%ADX/1=%22?%C4%90%F0%90%A9%87%26#%F0%91%B2%9E1%CC%88%CC%81D
	// http://G.BLoG/%E0%AD%8C~%F0%9F%AA%95%22%D6%93%E0%A9%AD%E1%B3%930%D0%8A/%C2%BC%E0%B4%BE3%F0%9E%8B%B0%F0%91%86%BD%C2%B2%E0%B3%A9%CC%80D/%7C+%F0%9F%81%9D+%5D%CC%81%CB%85/%CC%80/%E1%B0%BF/%CD%82K%E0%A5%A9%CC%81#%CC%82
	// https://1.47.4.5:11/+%3E%E2%9F%BCK//A%DB%97%F0%90%AC%BC$%F0%91%99%83%E2%82%A5%F0%9D%A8%A8%E0%BC%95%E0%B5%A8%3C%E0%BE%B0%F0%97%8D%91%E2%9E%8E%E0%B9%91$v%CC%80/%CC%94Z%E4%87%A4?%F0%91%91%9BC%C2%B9%E2%8A%A5%F0%91%91%8F1%E0%A0%BE%CB%BE%C3%9D%E1%B3%A8%E0%AB%A6%CC%81%CC%86&%E2%A4%88%F0%91%8A%A9%24B%F0%9D%8B%AA%CC%9A&&%CC%80%C2%A7%E4%92%9B&#%E0%AB%AE%F0%92%91%A2
	// http://G.hM/%CC%80%E0%A0%B1%CC%82%CC%80%F0%9E%A5%9F/:%21J%E2%9D%87#L%CC%82%E9%98%A6%22
	// http://1.1.4.6:2/%E3%B5%B8%C3%B2%DE%A6%E1%9E%B9/%C2%A7/?%E2%80%A60%CC%80%C3%B7&%CC%81%CC%A2%21%E0%AF%AB%CC%81%C3%A4&%F0%9F%AA%95%EA%99%B4%CC%80%E0%A5%AD&%F0%AD%B9%A9%F0%91%87%AE&&%E1%B7%93%CC%8B%E2%87%94%E2%90%8E%EA%A3%A5%E0%B5%9A%3D%E5%8E%A4%D9%AAB%F0%A5%8F%9A%3D%C2%A4%C3%AE%F0%91%84%AD%DC%8A%21%E2%82%8D3&%E1%81%8F%23#%CC%80%E0%BF%8B+$
}
