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

func ExampleIPv4() {
	gen := rapid.IPv4()

	for i := 0; i < 5; i++ {
		addr := gen.Example(i)
		fmt.Println(addr.String())
	}

	// Output:
	// 0.23.24.7
	// 100.146.0.0
	// 0.222.65.1
	// 1.49.104.14
	// 11.56.0.83
}

func ExampleIPv6() {
	gen := rapid.IPv6()

	for i := 0; i < 5; i++ {
		addr := gen.Example(i)
		fmt.Println(addr.String())
	}

	// Output:
	// 17:1807:e2c4:8210:7202:f4b2:a0e2:8dc
	// 6492:0:fa37:b00:b5c3:4e6:6a01:c802
	// de:4101:9f5:3:104:5dc:b600:905
	// 131:680e:97ff:d200:ae1:4d00:2300:103
	// b38:53:ff07:200:8c28:ee:ad00:1b
}
