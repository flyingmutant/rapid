// Copyright 2020 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"

	"pgregory.net/rapid"
)

func ExampleCustom() {
	gen := rapid.Custom(func(t *rapid.T) int {
		return rapid.Just(42).Draw(t, "answer").(int)
	})

	fmt.Println(gen.Example())
	// Output: 42
}
