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
	// https://r125pz05Rz-0j1d-11.AArP:17
	// http://L2.aBArTh:7/%25F0%259F%25AA%2595%2522%25D6%2593%25E0%25A9%25AD%25E1%25B3%25930%25D0%258A/%25C2%25BC%25E0%25B4%25BE3%25F0%259E%258B%25B0%25F0%2591%2586%25BD%25C2%25B2%25E0%25B3%25A9%25CC%2580D/%257C+%25F0%259F%2581%259D+%255D%25CC%2581%25CB%2585/%25CC%2580/%25E1%25B0%25BF/%25CD%2582K%25E0%25A5%25A9%25CC%2581
	// https://pH20DR11.aaA
	// http://h.AcCounTaNtS:4/%25F0%259E%25A5%259F/:%2521J%25E2%259D%2587
	// http://A.xN--s9bRJ9C:2
}
