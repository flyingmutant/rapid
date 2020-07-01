// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "reflect"

type value interface{}

type generatorImpl interface {
	String() string
	type_() reflect.Type
	value(t *T) value
}

type Generator struct {
	impl generatorImpl
	typ  reflect.Type
	str  string
}

func newGenerator(impl generatorImpl) *Generator {
	return &Generator{
		impl: impl,
		typ:  impl.type_(),
	}
}

func (g *Generator) String() string {
	if g.str == "" {
		g.str = g.impl.String()
	}

	return g.str
}

func (g *Generator) type_() reflect.Type {
	return g.typ
}

func (g *Generator) Draw(t *T, label string) interface{} {
	if t.log && t.tb != nil {
		t.tb.Helper()
	}
	return t.draw(g, label)
}

func (g *Generator) value(t *T) value {
	i := t.s.beginGroup(g.str, true)

	v := g.impl.value(t)
	u := reflect.TypeOf(v)
	assertf(v != nil, "%v has generated a nil value", g)
	assertf(u.AssignableTo(g.typ), "%v has generated a value of type %v which is not assignable to %v", g, u, g.typ)

	t.s.endGroup(i, false)

	return v
}

func (g *Generator) Example(seed ...int) interface{} {
	s := baseSeed()
	if len(seed) > 0 {
		s = uint64(seed[0])
	}

	v, n, err := example(g, newT(nil, newRandomBitStream(s, false), false))
	assertf(err == nil, "%v failed to generate an example in %v tries: %v", g, n, err)

	return v
}

func (g *Generator) Filter(fn interface{}) *Generator {
	return filter(g, fn)
}

func (g *Generator) Map(fn interface{}) *Generator {
	return map_(g, fn)
}

func example(g *Generator, t *T) (value, int, error) {
	for i := 1; ; i++ {
		r, err := recoverValue(g, t)
		if err == nil {
			return r, i, nil
		} else if i == exampleMaxTries {
			return nil, i, err
		}
	}
}

func recoverValue(g *Generator, t *T) (v value, err *testError) {
	defer func() { err = panicToError(recover(), 3) }()

	return g.value(t), nil
}
