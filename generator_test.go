// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"reflect"
	"testing"
)

type trivialGenImpl struct{}

func (trivialGenImpl) String() string          { return "" }
func (trivialGenImpl) type_() reflect.Type     { return uint64Type }
func (trivialGenImpl) value(s bitStream) Value { return s.drawBits(64) }

func BenchmarkTrivialGenImplValue(b *testing.B) {
	s := newRandomBitStream(baseSeed(), false)
	g := trivialGenImpl{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.value(s)
	}
}

func BenchmarkGenerator_Value(b *testing.B) {
	s := newRandomBitStream(baseSeed(), false)
	g := newGenerator(trivialGenImpl{})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.value(s)
	}
}
