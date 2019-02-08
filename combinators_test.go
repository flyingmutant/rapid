// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "testing"

func BenchmarkHeavyChain3(b *testing.B) {
	s := newRandomBitStream(prngSeed(), false)
	g := Custom(func(data Data) int { return data.Draw(Ints(), "").(int) }).
		Map(func(i int) (int, int) { return i, i << 13 }).
		Map(func(x int, y int) int { return x + y })
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.value(s)
	}
}
