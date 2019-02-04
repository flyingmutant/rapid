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
	byteKind    = "Bytes"
	intKind     = "Ints"
	int8Kind    = "Int8s"
	int16Kind   = "Int16s"
	int32Kind   = "Int32s"
	int64Kind   = "Int64s"
	uintKind    = "Uints"
	uint8Kind   = "Uint8s"
	uint16Kind  = "Uint16s"
	uint32Kind  = "Uint32s"
	uint64Kind  = "Uint64s"
	uintptrKind = "Uintptrs"

	uintptrSize = 32 << (^uintptr(0) >> 32 & 1)
	uintSize    = 32 << (^uint(0) >> 32 & 1)
	intSize     = uintSize

	minInt     = -1 << (intSize - 1)
	maxInt     = 1<<(intSize-1) - 1
	maxUint    = 1<<intSize - 1
	maxUintptr = 1<<(uint(uintptrSize)) - 1
)

var (
	intType   = reflect.TypeOf(int(0))
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

func Booleans() *Generator                 { return newGenerator(&boolGen{}) }
func (g *boolGen) String() string          { return "Booleans()" }
func (g *boolGen) type_() reflect.Type     { return reflect.TypeOf(false) }
func (g *boolGen) value(s bitStream) Value { return s.drawBits(1) == 1 }

func Bytes() *Generator    { return newIntegerGen(byteKind) }
func Ints() *Generator     { return newIntegerGen(intKind) }
func Int8s() *Generator    { return newIntegerGen(int8Kind) }
func Int16s() *Generator   { return newIntegerGen(int16Kind) }
func Int32s() *Generator   { return newIntegerGen(int32Kind) }
func Int64s() *Generator   { return newIntegerGen(int64Kind) }
func Uints() *Generator    { return newIntegerGen(uintKind) }
func Uint8s() *Generator   { return newIntegerGen(uint8Kind) }
func Uint16s() *Generator  { return newIntegerGen(uint16Kind) }
func Uint32s() *Generator  { return newIntegerGen(uint32Kind) }
func Uint64s() *Generator  { return newIntegerGen(uint64Kind) }
func Uintptrs() *Generator { return newIntegerGen(uintptrKind) }

func BytesMin(min byte) *Generator       { return newUintMinGen(byteKind, uint64(min)) }
func IntsMin(min int) *Generator         { return newIntMinGen(intKind, int64(min)) }
func Int8sMin(min int8) *Generator       { return newIntMinGen(int8Kind, int64(min)) }
func Int16sMin(min int16) *Generator     { return newIntMinGen(int16Kind, int64(min)) }
func Int32sMin(min int32) *Generator     { return newIntMinGen(int32Kind, int64(min)) }
func Int64sMin(min int64) *Generator     { return newIntMinGen(int64Kind, int64(min)) }
func UintsMin(min uint) *Generator       { return newUintMinGen(uintKind, uint64(min)) }
func Uint8sMin(min uint8) *Generator     { return newUintMinGen(uint8Kind, uint64(min)) }
func Uint16sMin(min uint16) *Generator   { return newUintMinGen(uint16Kind, uint64(min)) }
func Uint32sMin(min uint32) *Generator   { return newUintMinGen(uint32Kind, uint64(min)) }
func Uint64sMin(min uint64) *Generator   { return newUintMinGen(uint64Kind, uint64(min)) }
func UintptrsMin(min uintptr) *Generator { return newUintMinGen(uintptrKind, uint64(min)) }

func BytesMax(max byte) *Generator       { return newUintMaxGen(byteKind, uint64(max)) }
func IntsMax(max int) *Generator         { return newIntMaxGen(intKind, int64(max)) }
func Int8sMax(max int8) *Generator       { return newIntMaxGen(int8Kind, int64(max)) }
func Int16sMax(max int16) *Generator     { return newIntMaxGen(int16Kind, int64(max)) }
func Int32sMax(max int32) *Generator     { return newIntMaxGen(int32Kind, int64(max)) }
func Int64sMax(max int64) *Generator     { return newIntMaxGen(int64Kind, int64(max)) }
func UintsMax(max uint) *Generator       { return newUintMaxGen(uintKind, uint64(max)) }
func Uint8sMax(max uint8) *Generator     { return newUintMaxGen(uint8Kind, uint64(max)) }
func Uint16sMax(max uint16) *Generator   { return newUintMaxGen(uint16Kind, uint64(max)) }
func Uint32sMax(max uint32) *Generator   { return newUintMaxGen(uint32Kind, uint64(max)) }
func Uint64sMax(max uint64) *Generator   { return newUintMaxGen(uint64Kind, uint64(max)) }
func UintptrsMax(max uintptr) *Generator { return newUintMaxGen(uintptrKind, uint64(max)) }

func BytesRange(min byte, max byte) *Generator {
	return newUintRangeGen(byteKind, uint64(min), uint64(max))
}
func IntsRange(min int, max int) *Generator {
	return newIntRangeGen(intKind, int64(min), int64(max))
}
func Int8sRange(min int8, max int8) *Generator {
	return newIntRangeGen(int8Kind, int64(min), int64(max))
}
func Int16sRange(min int16, max int16) *Generator {
	return newIntRangeGen(int16Kind, int64(min), int64(max))
}
func Int32sRange(min int32, max int32) *Generator {
	return newIntRangeGen(int32Kind, int64(min), int64(max))
}
func Int64sRange(min int64, max int64) *Generator {
	return newIntRangeGen(int64Kind, int64(min), int64(max))
}
func UintsRange(min uint, max uint) *Generator {
	return newUintRangeGen(uintKind, uint64(min), uint64(max))
}
func Uint8sRange(min uint8, max uint8) *Generator {
	return newUintRangeGen(uint8Kind, uint64(min), uint64(max))
}
func Uint16sRange(min uint16, max uint16) *Generator {
	return newUintRangeGen(uint16Kind, uint64(min), uint64(max))
}
func Uint32sRange(min uint32, max uint32) *Generator {
	return newUintRangeGen(uint32Kind, uint64(min), uint64(max))
}
func Uint64sRange(min uint64, max uint64) *Generator {
	return newUintRangeGen(uint64Kind, uint64(min), uint64(max))
}
func UintptrsRange(min uintptr, max uintptr) *Generator {
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
			return fmt.Sprintf("%vRange(%v, %v)", g.kind, g.smin, g.smax)
		} else {
			return fmt.Sprintf("%vRange(%v, %v)", g.kind, g.umin, g.umax)
		}
	} else if g.hasMin {
		if g.signed {
			return fmt.Sprintf("%vMin(%v)", g.kind, g.smin)
		} else {
			return fmt.Sprintf("%vMin(%v)", g.kind, g.umin)
		}
	} else if g.hasMax {
		if g.signed {
			return fmt.Sprintf("%vMax(%v)", g.kind, g.smax)
		} else {
			return fmt.Sprintf("%vMax(%v)", g.kind, g.umax)
		}
	}

	return fmt.Sprintf("%v()", g.kind)
}

func (g *integerGen) type_() reflect.Type {
	return g.typ
}

func (g *integerGen) value(s bitStream) Value {
	var i int64
	var u uint64

	if g.signed {
		i = genIntRange(s, g.smin, g.smax, true)
	} else {
		u = genUintRange(s, g.umin, g.umax, true)
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
		return int64(i)
	case uintType:
		return uint(u)
	case uint8Type:
		return uint8(u)
	case uint16Type:
		return uint16(u)
	case uint32Type:
		return uint32(u)
	case uint64Type:
		return uint64(u)
	default:
		assertf(g.typ == uintptrType, "unhandled integer type %v", g.typ)
		return uintptr(u)
	}
}
