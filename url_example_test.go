// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

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
	// MV2zb0-S2j.trAveLcHAnnEL
	// Z.CU
	// r.ABBotT
	// r.AcCoUNTaNT
	// R.
}

func ExampleDomainOf() {
	gen := rapid.DomainOf(6, 5)

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}

	// Output:
	// Dg5G.
	// Z.CU
	// Bs.
	// AI.HkT
	// R.
}

func ExampleURL() {
	gen := rapid.URL()

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}

	// Output:
	// {https   U.aaA:4 V0%E2%90%9A%226%E0%BC%B0%F0%91%82%B0%F0%97%B3%80%F0%92%91%ADX/1=%22  false   }
	// {http   C.AarP:11 1%EF%BD%9F/%F0%9F%AA%95%22%D6%93%E0%A9%AD%E1%B3%930%D0%8A/%C2%BC%E0%B4%BE3%F0%9E%8B%B0%F0%91%86%BD%C2%B2%E0%B3%A9%CC%80D/%7C+%F0%9F%81%9D+%5D%CC%81%CB%85/%CC%80/%E1%B0%BF/%CD%82K%E0%A5%A9%CC%81  false   }
	// {https   Bs.:11   false   }
	// {http   MC0zJ.aCcOUNtAnT:2 J%E2%9D%87  false   }
	// {http   t.Xn--RvC1e0am3E:3 %CC%82/%E2%80%A60%CC%80%C3%B7/%CC%81%CC%A2%21%E0%AF%AB%CC%81%C3%A4/%F0%9F%AA%95%EA%99%B4%CC%80%E0%A5%AD/%F0%AD%B9%A9%F0%91%87%AE/%E1%B7%93%CC%8B%E2%87%94%E2%90%8E%EA%A3%A5%E0%B5%9A=%E5%8E%A4%D9%AAB%F0%A5%8F%9A=%C2%A4%C3%AE%F0%91%84%AD%DC%8A%21%E2%82%8D3/%E1%81%8F%23  false   }
}
