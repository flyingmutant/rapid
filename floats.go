// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"math/bits"
	"sort"
)

const (
	float32ExpBits  = 8
	float32ExpMax   = 1<<float32ExpBits - 1
	float32ExpBias  = 1<<(float32ExpBits-1) - 1
	float32MantBits = 23
	float32MantMax  = 1<<float32MantBits - 1

	float64ExpBits  = 11
	float64ExpMax   = 1<<float64ExpBits - 1
	float64ExpBias  = 1<<(float64ExpBits-1) - 1
	float64MantBits = 52
	float64MantMax  = 1<<float64MantBits - 1
)

var (
	float32ExpEnc [float32ExpMax + 1]uint32
	float32ExpDec [float32ExpMax + 1]uint32

	float64ExpEnc [float64ExpMax + 1]uint32
	float64ExpDec [float64ExpMax + 1]uint32
)

func init() {
	fillFloatTables(float32ExpEnc[:], float32ExpDec[:], float32ExpMax, float32ExpBias)
	fillFloatTables(float64ExpEnc[:], float64ExpDec[:], float64ExpMax, float64ExpBias)
}

func fillFloatTables(enc []uint32, dec []uint32, maxE uint32, bias uint32) {
	for e := uint32(0); e <= maxE; e++ {
		enc[e] = e
	}

	sort.Slice(enc, func(i, j int) bool { return floatExpKey(enc[i], maxE, bias) < floatExpKey(enc[j], maxE, bias) })

	for i, e := range enc {
		dec[e] = uint32(i)
	}
}

func floatExpKey(e uint32, maxE uint32, bias uint32) uint32 {
	if e == maxE {
		return math.MaxInt32
	}
	if e < bias {
		return math.MaxInt32/2 - e
	} else {
		return e
	}
}

func encodeExp32(e uint32) uint32 { return float32ExpDec[e] }
func decodeExp32(e uint32) uint32 { return float32ExpEnc[e] }
func encodeExp64(e uint32) uint32 { return float64ExpDec[e] }
func decodeExp64(e uint32) uint32 { return float64ExpEnc[e] }

func reverseBits(x uint64, n uint) uint64 {
	assert(n <= 64 && bits.Len64(x) <= int(n))
	return bits.Reverse64(x) >> (64 - n)
}

func transformMant(e uint32, bias uint32, bits_ uint, m uint64) uint64 {
	var fractBits uint
	if e <= bias {
		fractBits = bits_
	} else if e < bias+uint32(bits_) {
		fractBits = bits_ - uint(e-bias)
	} else {
		fractBits = 0
	}

	fractPart := m & (1<<fractBits - 1)
	m_ := (m ^ fractPart) | reverseBits(fractPart, fractBits)

	assert(bits.Len64(m_) <= int(bits_))
	return m_
}

func transformMant32(e uint32, m uint64) uint64 {
	return transformMant(e, float32ExpBias, float32MantBits, m)
}
func transformMant64(e uint32, m uint64) uint64 {
	return transformMant(e, float64ExpBias, float64MantBits, m)
}

func float32ToLex(f float32) (uint32, uint64) {
	u := math.Float32bits(f) & math.MaxInt32
	e := uint32(u >> float32MantBits)
	m := uint64(u & float32MantMax)
	return encodeExp32(e), transformMant32(e, m)
}

func float64ToLex(f float64) (uint32, uint64) {
	u := math.Float64bits(f) & math.MaxInt64
	e := uint32(u >> float64MantBits)
	m := uint64(u & float64MantMax)
	return encodeExp64(e), transformMant64(e, m)
}

func lexToFloat32(e uint32, m uint64) float32 {
	e = decodeExp32(e & float32ExpMax)
	m = transformMant32(e, m&float32MantMax)
	return math.Float32frombits(e<<float32MantBits | uint32(m))
}

func lexToFloat64(e uint32, m uint64) float64 {
	e = decodeExp64(e & float64ExpMax)
	m = transformMant64(e, m&float64MantMax)
	return math.Float64frombits(uint64(e)<<float64MantBits | m)
}
