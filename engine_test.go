// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"strings"
	"testing"
)

func TestPanicTraceback(t *testing.T) {
	t.Parallel()

	s := createRandomBitStream(t)
	g := Booleans().Filter(func(bool) bool { return false })

	_, err := recoverValue(g, s)
	if err == nil {
		t.Fatalf("no error from impossible filter")
	}

	lines := strings.Split(err.traceback, "\n")
	if !strings.HasSuffix(lines[0], "/rapid.satisfy") {
		t.Errorf("bad traceback from recoverValue():\n%v", err.traceback)
	}
}

func BenchmarkCheckOverhead(b *testing.B) {
	g := Uints()
	f := func(t *T) {
		g.Draw(t, "")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		checkTB(b, f)
	}
}
