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

func SliceOf(elem *Generator) *Generator {
	return SliceOfN(elem, -1, -1)
}

func SliceOfN(elem *Generator, minLen int, maxLen int) *Generator {
	assertValidRange(minLen, maxLen)

	return newGenerator(&sliceGen{
		typ:    reflect.SliceOf(elem.type_()),
		minLen: minLen,
		maxLen: maxLen,
		elem:   elem,
	})
}

func SliceOfDistinct(elem *Generator, keyFn interface{}) *Generator {
	return SliceOfNDistinct(elem, -1, -1, keyFn)
}

func SliceOfNDistinct(elem *Generator, minLen int, maxLen int, keyFn interface{}) *Generator {
	assertValidRange(minLen, maxLen)

	keyTyp := elem.type_()
	if keyFn != nil {
		t := reflect.TypeOf(keyFn)
		assertCallable(t, elem.type_(), "keyFn")
		keyTyp = t.Out(0)
	}
	assertf(keyTyp.Comparable(), "key type should be comparable (got %v)", keyTyp)

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
			return fmt.Sprintf("SliceOf(%v)", g.elem)
		} else {
			return fmt.Sprintf("SliceOfN(%v, minLen=%v, maxLen=%v)", g.elem, g.minLen, g.maxLen)
		}
	} else {
		key := ""
		if g.keyFn.IsValid() {
			key = fmt.Sprintf(", key=func(%v) %v", g.elem.type_(), g.keyTyp)
		}

		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("SliceOfDistinct(%v%v)", g.elem, key)
		} else {
			return fmt.Sprintf("SliceOfNDistinct(%v, minLen=%v, maxLen=%v%v)", g.elem, g.minLen, g.maxLen, key)
		}
	}
}

func (g *sliceGen) type_() reflect.Type {
	return g.typ
}

func (g *sliceGen) value(t *T) value {
	repeat := newRepeat(g.minLen, g.maxLen, -1)

	var seen reflect.Value
	if g.keyTyp != nil {
		seen = reflect.MakeMapWithSize(reflect.MapOf(g.keyTyp, emptyStructType), repeat.avg())
	}

	sl := reflect.MakeSlice(g.typ, 0, repeat.avg())
	for repeat.more(t.s, g.elem.String()) {
		e := reflect.ValueOf(g.elem.value(t))
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

func MapOf(key *Generator, val *Generator) *Generator {
	return MapOfN(key, val, -1, -1)
}

func MapOfN(key *Generator, val *Generator, minLen int, maxLen int) *Generator {
	assertValidRange(minLen, maxLen)
	assertf(key.type_().Comparable(), "key type should be comparable (got %v)", key.type_())

	return newGenerator(&mapGen{
		typ:    reflect.MapOf(key.type_(), val.type_()),
		minLen: minLen,
		maxLen: maxLen,
		key:    key,
		val:    val,
	})
}

func MapOfValues(val *Generator, keyFn interface{}) *Generator {
	return MapOfNValues(val, -1, -1, keyFn)
}

func MapOfNValues(val *Generator, minLen int, maxLen int, keyFn interface{}) *Generator {
	assertValidRange(minLen, maxLen)

	keyTyp := val.type_()
	if keyFn != nil {
		t := reflect.TypeOf(keyFn)
		assertCallable(t, val.type_(), "keyFn")
		keyTyp = t.Out(0)
	}
	assertf(keyTyp.Comparable(), "key type should be comparable (got %v)", keyTyp)

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
			return fmt.Sprintf("MapOf(%v, %v)", g.key, g.val)
		} else {
			return fmt.Sprintf("MapOfN(%v, %v, minLen=%v, maxLen=%v)", g.key, g.val, g.minLen, g.maxLen)
		}
	} else {
		key := ""
		if g.keyFn.IsValid() {
			key = fmt.Sprintf(", key=func(%v) %v", g.val.type_(), g.keyTyp)
		}

		if g.minLen < 0 && g.maxLen < 0 {
			return fmt.Sprintf("MapOfValues(%v%v)", g.val, key)
		} else {
			return fmt.Sprintf("MapOfNValues(%v, minLen=%v, maxLen=%v%v)", g.val, g.minLen, g.maxLen, key)
		}
	}
}

func (g *mapGen) type_() reflect.Type {
	return g.typ
}

func (g *mapGen) value(t *T) value {
	label := g.val.String()
	if g.key != nil {
		label = g.key.String() + "," + label
	}

	repeat := newRepeat(g.minLen, g.maxLen, -1)

	m := reflect.MakeMapWithSize(g.typ, repeat.avg())
	for repeat.more(t.s, label) {
		var k, v reflect.Value
		if g.keyTyp == nil {
			k = reflect.ValueOf(g.key.value(t))
			v = reflect.ValueOf(g.val.value(t))
		} else {
			v = reflect.ValueOf(g.val.value(t))
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

func ArrayOf(count int, elem *Generator) *Generator {
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
	return fmt.Sprintf("ArrayOf(%v, %v)", g.count, g.elem)
}

func (g *arrayGen) type_() reflect.Type {
	return g.typ
}

func (g *arrayGen) value(t *T) value {
	a := reflect.Indirect(reflect.New(g.typ))

	if g.count == 0 {
		t.s.drawBits(0)
	} else {
		for i := 0; i < g.count; i++ {
			e := reflect.ValueOf(g.elem.value(t))
			a.Index(i).Set(e)
		}
	}

	return a.Interface()
}
