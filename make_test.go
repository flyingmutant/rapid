// Copyright 2022 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"reflect"
	"testing"

	"pgregory.net/rapid"
)

type PrivateFields struct {
	private bool
}

func TestMake(t *testing.T) {
	// Private fields are ignored (and don't panic).
	rapid.Make[PrivateFields]().Example()
}

func TestMakeCustom(t *testing.T) {
	ex := rapid.MakeCustom[PrivateFields](rapid.MakeConfig{
		Types: map[reflect.Type]*rapid.Generator[any]{
			reflect.TypeOf(PrivateFields{}): rapid.Just(PrivateFields{private: true}).AsAny(),
		},
	}).Example()

	if !ex.private {
		t.Errorf(".private should be true. got: %#v", ex)
	}
}
