// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"math"
	"reflect"
)

const (
	byteKind    = "Byte"
	intKind     = "Int"
	int8Kind    = "Int8"
	int16Kind   = "Int16"
	int32Kind   = "Int32"
	int64Kind   = "Int64"
	uintKind    = "Uint"
	uint8Kind   = "Uint8"
	uint16Kind  = "Uint16"
	uint32Kind  = "Uint32"
	uint64Kind  = "Uint64"
	uintptrKind = "Uintptr"

	uintptrSize = 32 << (^uintptr(0) >> 32 & 1)
	uintSize    = 32 << (^uint(0) >> 32 & 1)
	intSize     = uintSize

	minInt     = -1 << (intSize - 1)
	maxInt     = 1<<(intSize-1) - 1
	maxUint    = 1<<intSize - 1
	maxUintptr = 1<<(uint(uintptrSize)) - 1
)

var (
	intType   = reflect.TypeOf(0)
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType    = reflect.TypeOf(uint(0))
	uint8Type   = reflect.TypeOf(uint8(0))
	uint16Type  = reflect.TypeOf(uint16(0))
	uint32Type  = reflect.TypeOf(uint32(0))
	uint64Type  = reflect.TypeOf(uint64(0))
	uintptrType = reflect.TypeOf(uintptr(0))

	integerKindToInfo = map[string]integerKindInfo{
		byteKind:    {typ: uint8Type, size: 1, umax: math.MaxUint8},
		intKind:     {typ: intType, signed: true, size: intSize / 8, smin: minInt, smax: maxInt},
		int8Kind:    {typ: int8Type, signed: true, size: 1, smin: math.MinInt8, smax: math.MaxInt8},
		int16Kind:   {typ: int16Type, signed: true, size: 2, smin: math.MinInt16, smax: math.MaxInt16},
		int32Kind:   {typ: int32Type, signed: true, size: 4, smin: math.MinInt32, smax: math.MaxInt32},
		int64Kind:   {typ: int64Type, signed: true, size: 8, smin: math.MinInt64, smax: math.MaxInt64},
		uintKind:    {typ: uintType, size: uintSize / 8, umax: maxUint},
		uint8Kind:   {typ: uint8Type, size: 1, umax: math.MaxUint8},
		uint16Kind:  {typ: uint16Type, size: 2, umax: math.MaxUint16},
		uint32Kind:  {typ: uint32Type, size: 4, umax: math.MaxUint32},
		uint64Kind:  {typ: uint64Type, size: 8, umax: math.MaxUint64},
		uintptrKind: {typ: uintptrType, size: uintptrSize / 8, umax: maxUintptr},
	}
)

type integerKindInfo struct {
	typ    reflect.Type
	signed bool
	size   int
	smin   int64
	smax   int64
	umax   uint64
}

type boolGen struct{}

func Bool() *Generator                 { return newGenerator(&boolGen{}) }
func (g *boolGen) String() string      { return "Bool()" }
func (g *boolGen) type_() reflect.Type { return reflect.TypeOf(false) }
func (g *boolGen) value(t *T) value    { return t.s.drawBits(1) == 1 }

func Byte() *Generator    { return newIntegerGen(byteKind) }
func Int() *Generator     { return newIntegerGen(intKind) }
func Int8() *Generator    { return newIntegerGen(int8Kind) }
func Int16() *Generator   { return newIntegerGen(int16Kind) }
func Int32() *Generator   { return newIntegerGen(int32Kind) }
func Int64() *Generator   { return newIntegerGen(int64Kind) }
func Uint() *Generator    { return newIntegerGen(uintKind) }
func Uint8() *Generator   { return newIntegerGen(uint8Kind) }
func Uint16() *Generator  { return newIntegerGen(uint16Kind) }
func Uint32() *Generator  { return newIntegerGen(uint32Kind) }
func Uint64() *Generator  { return newIntegerGen(uint64Kind) }
func Uintptr() *Generator { return newIntegerGen(uintptrKind) }

func ByteMin(min byte) *Generator       { return newUintMinGen(byteKind, uint64(min)) }
func IntMin(min int) *Generator         { return newIntMinGen(intKind, int64(min)) }
func Int8Min(min int8) *Generator       { return newIntMinGen(int8Kind, int64(min)) }
func Int16Min(min int16) *Generator     { return newIntMinGen(int16Kind, int64(min)) }
func Int32Min(min int32) *Generator     { return newIntMinGen(int32Kind, int64(min)) }
func Int64Min(min int64) *Generator     { return newIntMinGen(int64Kind, min) }
func UintMin(min uint) *Generator       { return newUintMinGen(uintKind, uint64(min)) }
func Uint8Min(min uint8) *Generator     { return newUintMinGen(uint8Kind, uint64(min)) }
func Uint16Min(min uint16) *Generator   { return newUintMinGen(uint16Kind, uint64(min)) }
func Uint32Min(min uint32) *Generator   { return newUintMinGen(uint32Kind, uint64(min)) }
func Uint64Min(min uint64) *Generator   { return newUintMinGen(uint64Kind, min) }
func UintptrMin(min uintptr) *Generator { return newUintMinGen(uintptrKind, uint64(min)) }

func ByteMax(max byte) *Generator       { return newUintMaxGen(byteKind, uint64(max)) }
func IntMax(max int) *Generator         { return newIntMaxGen(intKind, int64(max)) }
func Int8Max(max int8) *Generator       { return newIntMaxGen(int8Kind, int64(max)) }
func Int16Max(max int16) *Generator     { return newIntMaxGen(int16Kind, int64(max)) }
func Int32Max(max int32) *Generator     { return newIntMaxGen(int32Kind, int64(max)) }
func Int64Max(max int64) *Generator     { return newIntMaxGen(int64Kind, max) }
func UintMax(max uint) *Generator       { return newUintMaxGen(uintKind, uint64(max)) }
func Uint8Max(max uint8) *Generator     { return newUintMaxGen(uint8Kind, uint64(max)) }
func Uint16Max(max uint16) *Generator   { return newUintMaxGen(uint16Kind, uint64(max)) }
func Uint32Max(max uint32) *Generator   { return newUintMaxGen(uint32Kind, uint64(max)) }
func Uint64Max(max uint64) *Generator   { return newUintMaxGen(uint64Kind, max) }
func UintptrMax(max uintptr) *Generator { return newUintMaxGen(uintptrKind, uint64(max)) }

func ByteRange(min byte, max byte) *Generator {
	return newUintRangeGen(byteKind, uint64(min), uint64(max))
}
func IntRange(min int, max int) *Generator {
	return newIntRangeGen(intKind, int64(min), int64(max))
}
func Int8Range(min int8, max int8) *Generator {
	return newIntRangeGen(int8Kind, int64(min), int64(max))
}
func Int16Range(min int16, max int16) *Generator {
	return newIntRangeGen(int16Kind, int64(min), int64(max))
}
func Int32Range(min int32, max int32) *Generator {
	return newIntRangeGen(int32Kind, int64(min), int64(max))
}
func Int64Range(min int64, max int64) *Generator {
	return newIntRangeGen(int64Kind, min, max)
}
func UintRange(min uint, max uint) *Generator {
	return newUintRangeGen(uintKind, uint64(min), uint64(max))
}
func Uint8Range(min uint8, max uint8) *Generator {
	return newUintRangeGen(uint8Kind, uint64(min), uint64(max))
}
func Uint16Range(min uint16, max uint16) *Generator {
	return newUintRangeGen(uint16Kind, uint64(min), uint64(max))
}
func Uint32Range(min uint32, max uint32) *Generator {
	return newUintRangeGen(uint32Kind, uint64(min), uint64(max))
}
func Uint64Range(min uint64, max uint64) *Generator {
	return newUintRangeGen(uint64Kind, min, max)
}
func UintptrRange(min uintptr, max uintptr) *Generator {
	return newUintRangeGen(uintptrKind, uint64(min), uint64(max))
}

func newIntegerGen(kind string) *Generator {
	return newGenerator(&integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
	})
}

func newIntRangeGen(kind string, min int64, max int64) *Generator {
	assertf(min <= max, "invalid integer range [%v, %v]", min, max)

	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMin:          true,
		hasMax:          true,
	}
	g.smin = min
	g.smax = max

	return newGenerator(g)
}

func newIntMinGen(kind string, min int64) *Generator {
	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMin:          true,
	}
	g.smin = min

	return newGenerator(g)
}

func newIntMaxGen(kind string, max int64) *Generator {
	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMax:          true,
	}
	g.smax = max

	return newGenerator(g)
}

func newUintRangeGen(kind string, min uint64, max uint64) *Generator {
	assertf(min <= max, "invalid integer range [%v, %v]", min, max)

	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMin:          true,
		hasMax:          true,
	}
	g.umin = min
	g.umax = max

	return newGenerator(g)
}

func newUintMinGen(kind string, min uint64) *Generator {
	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMin:          true,
	}
	g.umin = min

	return newGenerator(g)
}

func newUintMaxGen(kind string, max uint64) *Generator {
	g := &integerGen{
		integerKindInfo: integerKindToInfo[kind],
		kind:            kind,
		hasMax:          true,
	}
	g.umax = max

	return newGenerator(g)
}

type integerGen struct {
	integerKindInfo
	kind   string
	umin   uint64
	hasMin bool
	hasMax bool
}

func (g *integerGen) String() string {
	if g.hasMin && g.hasMax {
		if g.signed {
			return fmt.Sprintf("%sRange(%d, %d)", g.kind, g.smin, g.smax)
		} else {
			return fmt.Sprintf("%sRange(%d, %d)", g.kind, g.umin, g.umax)
		}
	} else if g.hasMin {
		if g.signed {
			return fmt.Sprintf("%sMin(%d)", g.kind, g.smin)
		} else {
			return fmt.Sprintf("%sMin(%d)", g.kind, g.umin)
		}
	} else if g.hasMax {
		if g.signed {
			return fmt.Sprintf("%sMax(%d)", g.kind, g.smax)
		} else {
			return fmt.Sprintf("%sMax(%d)", g.kind, g.umax)
		}
	}

	return fmt.Sprintf("%s()", g.kind)
}

func (g *integerGen) type_() reflect.Type {
	return g.typ
}

func (g *integerGen) value(t *T) value {
	var i int64
	var u uint64

	if g.signed {
		i, _, _ = genIntRange(t.s, g.smin, g.smax, true)
	} else {
		u, _, _ = genUintRange(t.s, g.umin, g.umax, true)
	}

	switch g.typ {
	case intType:
		return int(i)
	case int8Type:
		return int8(i)
	case int16Type:
		return int16(i)
	case int32Type:
		return int32(i)
	case int64Type:
		return i
	case uintType:
		return uint(u)
	case uint8Type:
		return uint8(u)
	case uint16Type:
		return uint16(u)
	case uint32Type:
		return uint32(u)
	case uint64Type:
		return u
	default:
		assertf(g.typ == uintptrType, "unhandled integer type %v", g.typ)
		return uintptr(u)
	}
}
