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
	// https://r125pz05Rz-0j1d-11.AArP:17#0%E2%81%9B%F3%A0%84%9A
	// http://L2.aBArTh:7/%F0%9F%AA%95%22%D6%93%E0%A9%AD%E1%B3%930%D0%8A/%C2%BC%E0%B4%BE3%F0%9E%8B%B0%F0%91%86%BD%C2%B2%E0%B3%A9%CC%80D/%7C+%F0%9F%81%9D+%5D%CC%81%CB%85/%CC%80/%E1%B0%BF/%CD%82K%E0%A5%A9%CC%81#%CC%82
	// https://pH20DR11.aaA?%DB%97%F0%90%AC%BC%24%F0%91%99%83%E2%82%A5%F0%9D%A8%A8%E0%BC%95%E0%B5%A8%3C%E0%BE%B0%F0%97%8D%91%E2%9E%8E%E0%B9%91%24v%CC%80&%CC%94Z%E4%87%A4#%CC%A0%E1%81%AD
	// http://h.AcCounTaNtS:4/%F0%9E%A5%9F/:%21J%E2%9D%87#L%CC%82%E9%98%A6%22
	// http://A.xN--s9bRJ9C:2
}
