// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"reflect"
)

func SlicesOf(elem *Generator) *Generator {
	return SlicesOfN(elem, -1, -1)
}

func SlicesOfN(elem *Generator, minLen int, maxLen int) *Generator {
	assertValidRange(minLen, maxLen)

	return newGenerator(&sliceGen{
		typ:    reflect.SliceOf(elem.type_()),
		minLen: minLen,
		maxLen: maxLen,
		elem:   elem,
	})
}

func SlicesOfDistinct(elem *Generator, keyFn interface{}) *Generator {
	return SlicesOfNDistinct(elem, -1, -1, keyFn)
}

func SlicesOfNDistinct(elem *Generator, minLen int, maxLen int, keyFn interface{}) *Generator {
	assertValidRange(minLen, maxLen)

	keyTyp := elem.type_()
	if keyFn != nil {
		t := reflect.TypeOf(keyFn)
		assertCallable(t, elem.type_(), "keyFn", 0)
		assertf(t.NumOut() == 1, "keyFn should have 1 output parameter, not %v", t.NumOut())
		keyTyp = t.Out(0)
	}
	assertf(keyTyp.Comparable(), "key type should be comparable (got %v)", prettyType{keyTyp})

	return newGenerator(&sliceGen{
		typ:    reflect.SliceOf(elem.type_()),
		minLen: minLen,
		maxLen: maxLen,
		elem:   elem,
		keyTyp: keyTyp,
		keyFn:  reflect.ValueOf(keyFn),
	})
}

type sliceGen struct {
	typ    reflect.Type
	minLen int
	maxLen int
	elem   *Generator
	keyTyp reflect.Type
	keyFn  reflect.Value
}

func (g *sliceGen) String() string {
	if g.keyTyp == nil {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("SlicesOf(%v)", g.elem)
		} else {
			return fmt.Sprintf("SlicesOfN(%v, minLen=%v, maxLen=%v)", g.elem, g.minLen, g.maxLen)
		}
	} else {
		key := ""
		if g.keyFn.IsValid() {
			key = fmt.Sprintf(", key=func(%v) %v", prettyType{g.elem.type_()}, prettyType{g.keyTyp})
		}

		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("SlicesOfDistinct(%v%v)", g.elem, key)
		} else {
			return fmt.Sprintf("SlicesOfNDistinct(%v, minLen=%v, maxLen=%v%v)", g.elem, g.minLen, g.maxLen, key)
		}
	}
}

func (g *sliceGen) type_() reflect.Type {
	return g.typ
}

func (g *sliceGen) value(s bitStream) Value {
	repeat := newRepeat(g.minLen, g.maxLen, -1)

	var seen reflect.Value
	if g.keyTyp != nil {
		seen = reflect.MakeMapWithSize(reflect.MapOf(g.keyTyp, emptyStructType), repeat.avg())
	}

	sl := reflect.MakeSlice(g.typ, 0, repeat.avg())
	for repeat.more(s, g.elem.String()) {
		e := reflect.ValueOf(g.elem.value(s))
		if g.keyTyp == nil {
			sl = reflect.Append(sl, e)
		} else {
			k := e
			if g.keyFn.IsValid() {
				k = g.keyFn.Call([]reflect.Value{k})[0]
			}

			if seen.MapIndex(k).IsValid() {
				repeat.reject()
			} else {
				seen.SetMapIndex(k, emptyStructValue)
				sl = reflect.Append(sl, e)
			}
		}
	}

	return sl.Interface()
}

func MapsOf(key *Generator, val *Generator) *Generator {
	return MapsOfN(key, val, -1, -1)
}

func MapsOfN(key *Generator, val *Generator, minLen int, maxLen int) *Generator {
	assertValidRange(minLen, maxLen)
	assertf(key.type_().Comparable(), "key type should be comparable (got %v)", prettyType{key.type_()})

	return newGenerator(&mapGen{
		typ:    reflect.MapOf(key.type_(), val.type_()),
		minLen: minLen,
		maxLen: maxLen,
		key:    key,
		val:    val,
	})
}

func MapsOfValues(val *Generator, keyFn interface{}) *Generator {
	return MapsOfNValues(val, -1, -1, keyFn)
}

func MapsOfNValues(val *Generator, minLen int, maxLen int, keyFn interface{}) *Generator {
	assertValidRange(minLen, maxLen)

	keyTyp := val.type_()
	if keyFn != nil {
		t := reflect.TypeOf(keyFn)
		assertCallable(t, val.type_(), "keyFn", 0)
		assertf(t.NumOut() == 1, "keyFn should have 1 output parameter, not %v", t.NumOut())
		keyTyp = t.Out(0)
	}
	assertf(keyTyp.Comparable(), "key type should be comparable (got %v)", prettyType{keyTyp})

	return newGenerator(&mapGen{
		typ:    reflect.MapOf(keyTyp, val.type_()),
		minLen: minLen,
		maxLen: maxLen,
		val:    val,
		keyTyp: keyTyp,
		keyFn:  reflect.ValueOf(keyFn),
	})
}

type mapGen struct {
	typ    reflect.Type
	minLen int
	maxLen int
	key    *Generator
	val    *Generator
	keyTyp reflect.Type
	keyFn  reflect.Value
}

func (g *mapGen) String() string {
	if g.keyTyp == nil {
		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("MapsOf(%v, %v)", g.key, g.val)
		} else {
			return fmt.Sprintf("MapsOfN(%v, %v, minLen=%v, maxLen=%v)", g.key, g.val, g.minLen, g.maxLen)
		}
	} else {
		key := ""
		if g.keyFn.IsValid() {
			key = fmt.Sprintf(", key=func(%v) %v", prettyType{g.val.type_()}, prettyType{g.keyTyp})
		}

		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("MapsOfValues(%v%v)", g.val, key)
		} else {
			return fmt.Sprintf("MapsOfNValues(%v, minLen=%v, maxLen=%v%v)", g.val, g.minLen, g.maxLen, key)
		}
	}
}

func (g *mapGen) type_() reflect.Type {
	return g.typ
}

func (g *mapGen) value(s bitStream) Value {
	label := g.val.String()
	if g.key != nil {
		label = g.key.String() + "," + label
	}

	repeat := newRepeat(g.minLen, g.maxLen, -1)

	m := reflect.MakeMapWithSize(g.typ, repeat.avg())
	for repeat.more(s, label) {
		var k, v reflect.Value
		if g.keyTyp == nil {
			k = reflect.ValueOf(g.key.value(s))
			v = reflect.ValueOf(g.val.value(s))
		} else {
			v = reflect.ValueOf(g.val.value(s))
			k = v
			if g.keyFn.IsValid() {
				k = g.keyFn.Call([]reflect.Value{v})[0]
			}
		}

		if m.MapIndex(k).IsValid() {
			repeat.reject()
		} else {
			m.SetMapIndex(k, v)
		}
	}

	return m.Interface()
}

func ArraysOf(count int, elem *Generator) *Generator {
	assertf(count >= 0 && count < 1024, "array element count should be in [0, 1024] (got %v)", count)

	return newGenerator(&arrayGen{
		typ:   reflect.ArrayOf(count, elem.type_()),
		count: count,
		elem:  elem,
	})
}

type arrayGen struct {
	typ   reflect.Type
	count int
	elem  *Generator
}

func (g *arrayGen) String() string {
	return fmt.Sprintf("ArraysOf(%v, %v)", g.count, g.elem)
}

func (g *arrayGen) type_() reflect.Type {
	return g.typ
}

func (g *arrayGen) value(s bitStream) Value {
	a := reflect.Indirect(reflect.New(g.typ))

	if g.count == 0 {
		s.drawBits(0)
	} else {
		for i := 0; i < g.count; i++ {
			e := reflect.ValueOf(g.elem.value(s))
			a.Index(i).Set(e)
		}
	}

	return a.Interface()
}
