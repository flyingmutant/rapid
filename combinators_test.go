// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "testing"

type intPair struct {
	x int
	y int
}

func BenchmarkHeavyChain3(b *testing.B) {
	t := newT(nil, newRandomBitStream(baseSeed(), false), false)
	g := Custom(func(t *T) int { return Int().Draw(t, "").(int) }).
		Map(func(i int) intPair { return intPair{i, i << 13} }).
		Map(func(p intPair) int { return p.x + p.y })
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.value(t)
	}
}
