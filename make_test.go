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

type privateFields struct {
	private bool
}

type fieldOverride struct {
	Name  string
	Count int
}

type stringAlias string

type kindOverride struct {
	Name  string
	Alias stringAlias
}

func TestMakeIgnoresPrivateFields(t *testing.T) {
	// Private fields are ignored (and don't panic).
	rapid.Make[privateFields]().Example()
}

func TestMakeCustomTypeOverride(t *testing.T) {
	ex := rapid.MakeCustom[privateFields](rapid.MakeConfig{
		Types: map[reflect.Type]*rapid.Generator[any]{
			reflect.TypeOf(privateFields{}): rapid.Just(privateFields{private: true}).AsAny(),
		},
	}).Example()

	if !ex.private {
		t.Errorf(".private should be true. got: %#v", ex)
	}
}

func TestMakeCustomFieldOverride(t *testing.T) {
	const customName = "custom-name"

	v := rapid.MakeCustom[fieldOverride](rapid.MakeConfig{
		Fields: map[reflect.Type]map[string]*rapid.Generator[any]{
			reflect.TypeOf(fieldOverride{}): {
				"Name": rapid.Just(customName).AsAny(),
			},
		},
	}).Example()

	if v.Name != customName {
		t.Fatalf("expected Name to be %q, got %q", customName, v.Name)
	}
}

func TestMakeCustomKindOverride(t *testing.T) {
	const customValue = "kind-override"

	v := rapid.MakeCustom[kindOverride](rapid.MakeConfig{
		Kinds: map[reflect.Kind]*rapid.Generator[any]{
			reflect.String: rapid.Just(customValue).AsAny(),
		},
	}).Example()

	if v.Name != customValue {
		t.Fatalf("expected Name to be %q, got %q", customValue, v.Name)
	}

	expectedAlias := stringAlias(customValue)
	if v.Alias != expectedAlias {
		t.Fatalf("expected Alias to be %q, got %q", expectedAlias, v.Alias)
	}
}
