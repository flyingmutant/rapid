// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "testing"

func TestLexToFloat32Roundtrip(t *testing.T) {
	Check(t, func(t *T, sign bool, e uint32, m uint64) {
		f := lexToFloat32(sign, e, m)
		sign_, e_, m_ := float32ToLex(f)
		if sign_ != sign || e_ != e || m_ != m {
			t.Fatalf("lex encoding roundtrip failed: (%v, %v, %v) vs expected (%v, %v, %v)", sign_, e_, m_, sign, e, m)
		}
	}, Booleans(), Uint32sMax(float32ExpMask), Uint64sMax(float32MantMask))
}

func TestLexToFloat64Roundtrip(t *testing.T) {
	Check(t, func(t *T, sign bool, e uint32, m uint64) {
		f := lexToFloat64(sign, e, m)
		sign_, e_, m_ := float64ToLex(f)
		if sign_ != sign || e_ != e || m_ != m {
			t.Fatalf("lex encoding roundtrip failed: (%v, %v, %v) vs expected (%v, %v, %v)", sign_, e_, m_, sign, e, m)
		}
	}, Booleans(), Uint32sMax(float64ExpMask), Uint64sMax(float64MantMask))
}
