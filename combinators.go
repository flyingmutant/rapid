// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"reflect"
	"strings"
)

const tryLabel = "try"

var (
	boolType = reflect.TypeOf(false)
	tPtrType = reflect.TypeOf((*T)(nil))
)

func Custom(fn interface{}) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, tPtrType, "fn")

	return newGenerator(&customGen{
		typ: t.Out(0),
		fn:  f,
	})
}

type customGen struct {
	typ reflect.Type
	fn  reflect.Value
}

func (g *customGen) String() string {
	return fmt.Sprintf("Custom(%v)", g.typ)
}

func (g *customGen) type_() reflect.Type {
	return g.typ
}

func (g *customGen) value(s bitStream) Value {
	return find(g.maybeValue, s, small)
}

func (g *customGen) maybeValue(s bitStream) Value {
	t := newT(nil, s, false)

	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(invalidData); !ok {
				panic(r)
			}
		}
	}()

	return call(g.fn, reflect.ValueOf(t))
}

func filter(g *Generator, fn interface{}) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, g.type_(), "fn")
	assertf(t.Out(0) == boolType, "fn should return bool, not %v", t.Out(0))

	return newGenerator(&filteredGen{
		g: g,
		fn: func(v Value) bool {
			return call(f, reflect.ValueOf(v)).(bool)
		},
	})
}

type filteredGen struct {
	g  *Generator
	fn func(Value) bool
}

func (g *filteredGen) String() string {
	return fmt.Sprintf("%v.Filter(...)", g.g)
}

func (g *filteredGen) type_() reflect.Type {
	return g.g.type_()
}

func (g *filteredGen) value(s bitStream) Value {
	return find(g.maybeValue, s, small)
}

func (g *filteredGen) maybeValue(s bitStream) Value {
	v := g.g.value(s)
	if g.fn(v) {
		return v
	} else {
		return nil
	}
}

func find(gen func(bitStream) Value, s bitStream, tries int) Value {
	for n := 0; n < tries; n++ {
		i := s.beginGroup(tryLabel, false)
		v := gen(s)
		ok := v != nil
		s.endGroup(i, !ok)

		if ok {
			return v
		}
	}

	panic(invalidData(fmt.Sprintf("failed to find suitable value in %d tries", tries)))
}

func map_(g *Generator, fn interface{}) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, g.type_(), "fn")

	return newGenerator(&mappedGen{
		typ: t.Out(0),
		g:   g,
		fn:  f,
	})
}

type mappedGen struct {
	typ reflect.Type
	g   *Generator
	fn  reflect.Value
}

func (g *mappedGen) String() string {
	return fmt.Sprintf("%v.Map(func(...) %v)", g.g, g.typ)
}

func (g *mappedGen) type_() reflect.Type {
	return g.typ
}

func (g *mappedGen) value(s bitStream) Value {
	v := reflect.ValueOf(g.g.value(s))
	return call(g.fn, v)
}

func Just(value Value) *Generator {
	return newGenerator(&sampledGen{
		typ:    reflect.TypeOf(value),
		values: []Value{value},
	})
}

func SampledFrom(slice interface{}) *Generator {
	v := reflect.ValueOf(slice)
	t := v.Type()

	assertf(t.Kind() == reflect.Slice, "argument should be a slice, not %v", t.Kind())
	assertf(v.Len() > 0, "slice should not be empty")

	values := make([]Value, v.Len())
	for i := 0; i < v.Len(); i++ {
		values[i] = v.Index(i).Interface()
	}

	return newGenerator(&sampledGen{
		typ:    t.Elem(),
		values: values,
	})
}

type sampledGen struct {
	typ    reflect.Type
	values []Value
}

func (g *sampledGen) String() string {
	if len(g.values) == 1 {
		return fmt.Sprintf("Just(%v)", g.values[0])
	} else {
		return fmt.Sprintf("SampledFrom(%v %v)", len(g.values), g.typ)
	}
}

func (g *sampledGen) type_() reflect.Type {
	return g.typ
}

func (g *sampledGen) value(s bitStream) Value {
	i := genIndex(s, len(g.values), true)

	return g.values[i]
}

func OneOf(gens ...*Generator) *Generator {
	assertf(len(gens) > 0, "at least one generator should be specified")
	for i, g := range gens {
		assertf(g.type_() == gens[0].type_(), "generator %v (%v) should generate %v, not %v", i, g, gens[0].type_(), g.type_())
	}

	return newGenerator(&oneOfGen{
		typ:  gens[0].type_(),
		gens: gens,
	})
}

type oneOfGen struct {
	typ  reflect.Type
	gens []*Generator
}

func (g *oneOfGen) String() string {
	strs := make([]string, len(g.gens))
	for i, g := range g.gens {
		strs[i] = g.String()
	}

	return fmt.Sprintf("OneOf(%v)", strings.Join(strs, ", "))
}

func (g *oneOfGen) type_() reflect.Type {
	return g.typ
}

func (g *oneOfGen) value(s bitStream) Value {
	i := genIndex(s, len(g.gens), true)

	return g.gens[i].value(s)
}

func Ptrs(elem *Generator, allowNil bool) *Generator {
	return newGenerator(&ptrGen{
		typ:      reflect.PtrTo(elem.type_()),
		elem:     elem,
		allowNil: allowNil,
	})
}

type ptrGen struct {
	typ      reflect.Type
	elem     *Generator
	allowNil bool
}

func (g *ptrGen) String() string {
	return fmt.Sprintf("Ptrs(%v, allowNil=%v)", g.elem, g.allowNil)
}

func (g *ptrGen) type_() reflect.Type {
	return g.typ
}

func (g *ptrGen) value(s bitStream) Value {
	pNonNil := float64(1)
	if g.allowNil {
		pNonNil = 0.5
	}

	if flipBiasedCoin(s, pNonNil) {
		p := reflect.New(g.elem.type_())
		p.Elem().Set(reflect.ValueOf(g.elem.value(s)))
		return p.Interface()
	} else {
		return reflect.Zero(g.typ).Interface()
	}
}
