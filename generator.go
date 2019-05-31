// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "reflect"

type Value interface{}

type generatorImpl interface {
	String() string
	type_() reflect.Type
	value(bitStream) Value
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
		str:  impl.String(),
	}
}

func (g *Generator) String() string {
	return g.str
}

func (g *Generator) type_() reflect.Type {
	return g.typ
}

func (g *Generator) Draw(src Source, label string, unpack ...interface{}) Value {
	return src.Draw(g, label, unpack...)
}

func (g *Generator) value(s bitStream) Value {
	i := s.beginGroup(g.str, true)

	v := g.impl.value(s)
	t := reflect.TypeOf(v)
	assertf(v != nil, "%v has generated a nil value", g)
	assertf(t.AssignableTo(g.typ), "%v has generated a value of type %v which is not assignable to %v", g, prettyType{t}, prettyType{g.typ})

	s.endGroup(i, false)

	return v
}

func (g *Generator) Example(seed ...uint64) (Value, int, error) {
	s := prngSeed()
	if len(seed) > 0 {
		s = seed[0]
	}

	return example(g, newRandomBitStream(s, false))
}

func (g *Generator) Filter(fn interface{}) *Generator {
	return filter(g, fn, small, "")
}

func (g *Generator) Map(fn interface{}) *Generator {
	return map_(g, fn)
}

func example(g *Generator, s bitStream) (Value, int, error) {
	for i := 1; ; i++ {
		r, err := recoverValue(g, s)
		if err == nil {
			return r, i, nil
		} else if i == exampleMaxTries {
			if err != nil {
				return nil, i, err
			}
			return nil, i, errCantGenDueToFilter
		}
	}
}

func recoverValue(g *Generator, s bitStream) (v Value, err *testError) {
	defer func() { err = panicToError(recover(), 3) }()

	return g.value(s), nil
}
