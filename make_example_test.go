// Copyright 2022 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"

	"pgregory.net/rapid"
)

func ExampleMake() {
	gen := rapid.Make[map[int32]bool]()

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// map[-53:false -49:true -23:false -5:false 1:true 56:false]
	// map[-3:true 0:true]
	// map[0:true]
	// map[-103:true -26:true -7:true -1:false 78:false 22525032:true]
	// map[]
}

type nodeValue int32

type tree struct {
	Value       nodeValue
	Left, Right *tree
}

func (t *tree) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("(%s %v %s)", t.Left.String(), t.Value, t.Right.String())
}

func ExampleMake_tree() {
	gen := rapid.Make[*tree]()

	for i := 0; i < 5; i++ {
		fmt.Println(gen.Example(i))
	}
	// Output:
	// (nil 1 (nil 56 nil))
	// (((nil -1 (((((nil -101 ((nil -2 ((((nil -1 nil) -802741815 nil) -106379 ((nil 1288 nil) 17 (((((nil -58 nil) -55219410 nil) 7713 nil) 238 (((nil -2 nil) 47243076 nil) 57485 nil)) -242 nil))) -29499386 nil)) -5 nil)) 926545 nil) 14314 (nil 3 nil)) -12 (nil -2 ((nil 1 nil) -2 (((nil 3 nil) -7083 ((nil -70 (nil -2799 nil)) 1787 (nil -15025 nil))) 508 (nil 6 nil))))) 41418 nil)) -2 nil) -3 nil)
	// nil
	// (((nil -15 (nil 454 nil)) -3 (nil 2096073436 ((nil 6 ((nil -2147483648 nil) 1 nil)) 332 nil))) 78 nil)
	// nil
}
