// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "testing"

func TestExp32Roundtrip(t *testing.T) {
	Check(t, func(t *T, e uint32) {
		enc := encodeExp32(e)
		dec := decodeExp32(enc)
		if dec != e {
			t.Fatalf("exponent encoding roundtrip failed: %v vs expected %v", dec, e)
		}
	}, Uint32sMax(float32ExpMax))
}

func TestExp64Roundtrip(t *testing.T) {
	Check(t, func(t *T, e uint32) {
		enc := encodeExp64(e)
		dec := decodeExp64(enc)
		if dec != e {
			t.Fatalf("exponent encoding roundtrip failed: %v vs expected %v", dec, e)
		}
	}, Uint32sMax(float64ExpMax))
}

func TestLexToFloat32Roundtrip(t *testing.T) {
	Check(t, func(t *T, sign bool, e uint32, m uint64) {
		f := lexToFloat32(sign, e, m)
		sign_, e_, m_ := float32ToLex(f)
		if sign_ != sign || e_ != e || m_ != m {
			t.Fatalf("lex encoding roundtrip failed: (%v, %v, %v) vs expected (%v, %v, %v)", sign_, e_, m_, sign, e, m)
		}
	}, Booleans(), Uint32sMax(float32ExpMax), Uint64sMax(float32MantMax))
}

func TestLexToFloat64Roundtrip(t *testing.T) {
	Check(t, func(t *T, sign bool, e uint32, m uint64) {
		f := lexToFloat64(sign, e, m)
		sign_, e_, m_ := float64ToLex(f)
		if sign_ != sign || e_ != e || m_ != m {
			t.Fatalf("lex encoding roundtrip failed: (%v, %v, %v) vs expected (%v, %v, %v)", sign_, e_, m_, sign, e, m)
		}
	}, Booleans(), Uint32sMax(float64ExpMax), Uint64sMax(float64MantMax))
}
