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
	dataType = reflect.TypeOf((*Data)(nil)).Elem()
)

func Custom(fn interface{}) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, dataType, "fn", 0)
	assertf(t.NumOut() > 0, "fn should have at least one output parameter")

	return newGenerator(&customGen{
		typ: retTypeOf(t),
		fn:  f,
	})
}

type customGen struct {
	typ reflect.Type
	fn  reflect.Value
}

func (g *customGen) String() string {
	return fmt.Sprintf("Custom(func(...) %v)", prettyType{g.typ})
}

func (g *customGen) type_() reflect.Type {
	return g.typ
}

func (g *customGen) value(s bitStream) Value {
	data := &bitStreamData{s}

	return call(g.fn, reflect.ValueOf(data), g.typ)
}

func Tuple(gens ...*Generator) *Generator {
	assertf(len(gens) > 0, "at least one generator should be specified")

	if len(gens) == 1 && isTuple(gens[0].type_()) {
		return gens[0]
	}

	genTypes := make([]reflect.Type, len(gens))
	for i, g := range gens {
		genTypes[i] = g.type_()
	}

	return newGenerator(&tupleGen{
		typ:  tupleOf(genTypes),
		gens: gens,
	})
}

type tupleGen struct {
	typ  reflect.Type
	gens []*Generator
}

func (g *tupleGen) String() string {
	b := &strings.Builder{}

	b.WriteString("Tuple(")
	for i, f := range g.gens {
		b.WriteString(f.String())
		if i != len(g.gens)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(")")

	return b.String()
}

func (g *tupleGen) type_() reflect.Type {
	return g.typ
}

func (g *tupleGen) value(s bitStream) Value {
	v := reflect.Indirect(reflect.New(g.typ))

	for i, g := range g.gens {
		v.Field(i).Set(reflect.ValueOf(g.value(s)))
	}

	return v.Interface()
}

func filter(g *Generator, fn interface{}, tries int, stopMsg string) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, g.type_(), "fn", 0)
	assertf(t.NumOut() == 1, "fn should have 1 output parameter, not %v", t.NumOut())
	assertf(t.Out(0) == boolType, "fn should return bool, not %v", t.Out(0))

	return newGenerator(&filteredGen{
		g: g,
		fn: func(v Value) bool {
			return call(f, reflect.ValueOf(v), nil).(bool)
		},
		tries:   tries,
		stopMsg: stopMsg,
	})
}

type filteredGen struct {
	g       *Generator
	fn      func(Value) bool
	tries   int
	stopMsg string
}

func (g *filteredGen) String() string {
	return fmt.Sprintf("%v.Filter(...)", g.g)
}

func (g *filteredGen) type_() reflect.Type {
	return g.g.type_()
}

func (g *filteredGen) value(s bitStream) Value {
	return satisfy(g.fn, g.g.value, s, g.tries, g.stopMsg)
}

func satisfy(filter func(Value) bool, gen func(bitStream) Value, s bitStream, tries int, stopMsg string) Value {
	for n := 0; n < tries; n++ {
		i := s.beginGroup(tryLabel, false)
		v := gen(s)
		ok := filter(v)
		s.endGroup(i, stopMsg == "" && !ok)

		if ok {
			return v
		}
	}

	if stopMsg != "" {
		panic(stopTest(stopMsg))
	} else {
		panic(invalidData(fmt.Sprintf("failed to satisfy filter in %d tries", tries)))
	}
}

func map_(g *Generator, fn interface{}) *Generator {
	f := reflect.ValueOf(fn)
	t := f.Type()

	assertCallable(t, g.type_(), "fn", 0)
	assertf(t.NumOut() > 0, "fn should have at least one output parameter")

	return newGenerator(&mappedGen{
		typ: retTypeOf(t),
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
	return fmt.Sprintf("%v.Map(func(...) %v)", g.g, prettyType{g.typ})
}

func (g *mappedGen) type_() reflect.Type {
	return g.typ
}

func (g *mappedGen) value(s bitStream) Value {
	v := reflect.ValueOf(g.g.value(s))
	return call(g.fn, v, g.typ)
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
		assertf(g.type_() == gens[0].type_(), "generator %v (%v) should generate %v, not %v", i, g, prettyType{gens[0].type_()}, prettyType{g.type_()})
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
