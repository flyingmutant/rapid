// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"strings"
	"testing"
)

func TestPanicTraceback(t *testing.T) {
	s := createRandomBitStream(t)
	g := Booleans().Filter(func(bool) bool { return false })

	_, err := recoverValue(g, s)
	tr := err.traceback()
	lines := strings.Split(tr, "\n")
	if !strings.HasSuffix(lines[0], "/rapid.(*filteredGen).value") {
		t.Errorf("bad traceback from recoverValue():\n%v", tr)
	}
}

func BenchmarkCheckOverhead(b *testing.B) {
	g := Uints()
	f := func(t *T) {
		t.Draw(g, "")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		checkTB(b, f)
	}
}
