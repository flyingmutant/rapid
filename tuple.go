// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tupleFieldPrefix    = "RapidTupleField"
	tupleFirstFieldName = tupleFieldPrefix + "0"
)

func isTuple(t reflect.Type) bool {
	return t.Kind() == reflect.Struct && t.NumField() > 0 && t.Field(0).Name == tupleFirstFieldName
}

func tupleOf(fieldTypes []reflect.Type) reflect.Type {
	fields := make([]reflect.StructField, len(fieldTypes))

	for i, t := range fieldTypes {
		fields[i] = reflect.StructField{
			Name: fmt.Sprintf("%s%d", tupleFieldPrefix, i),
			Type: t,
		}
	}

	return reflect.StructOf(fields)
}

func retTypeOf(fn reflect.Type) reflect.Type {
	assert(fn.Kind() == reflect.Func && fn.NumOut() > 0)

	if fn.NumOut() == 1 {
		return fn.Out(0)
	}

	fieldTypes := make([]reflect.Type, fn.NumOut())
	for i := 0; i < fn.NumOut(); i++ {
		fieldTypes[i] = fn.Out(i)
	}

	return tupleOf(fieldTypes)
}

func assertCallable(fn reflect.Type, t reflect.Type, name string, start int) {
	assertf(fn.Kind() == reflect.Func, "%v should be a function, not %v", name, fn.Kind())

	if isTuple(t) {
		assertf(fn.NumIn() == t.NumField()+start, "%v should have %v parameters, not %v", name, t.NumField()+start, fn.NumIn())
		for i := start; i < fn.NumIn(); i++ {
			assertf(t.Field(i-start).Type.AssignableTo(fn.In(i)), "parameter %v (%v) of %v should be assignable from %v", i, fn.In(i), name, t.Field(i-start).Type)
		}
	} else {
		assertf(fn.NumIn() == 1+start, "%v should have %v parameters, not %v", name, 1+start, fn.NumIn())
		assertf(t.AssignableTo(fn.In(start)), "parameter %v (%v) of %v should be assignable from %v", start, fn.In(start), name, t)
	}
}

func call(fn reflect.Value, arg reflect.Value, tuple reflect.Type) Value {
	t := arg.Type()

	var r []reflect.Value
	if isTuple(t) {
		n := t.NumField()
		args := make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			args[i] = arg.Field(i)
		}
		r = fn.Call(args)
	} else {
		r = fn.Call([]reflect.Value{arg})
	}

	if len(r) == 0 {
		return nil
	} else if len(r) == 1 {
		return r[0].Interface()
	} else {
		return packTuple(tuple, r...).Interface()
	}
}

func packTuple(tuple reflect.Type, fields ...reflect.Value) reflect.Value {
	tup := reflect.Indirect(reflect.New(tuple))
	for i, field := range fields {
		tup.Field(i).Set(field)
	}
	return tup
}

func unpackTuple(v reflect.Value, unpack ...interface{}) {
	t := v.Type()

	assertf(isTuple(t), "trying to unpack value of type %v which is not multi-valued into %v values", prettyType{t}, len(unpack))
	assertf(t.NumField() == len(unpack), "trying to unpack %v values of %v into %v", t.NumField(), prettyType{t}, len(unpack))

	for i, u := range unpack {
		ut := reflect.TypeOf(u)
		assertf(ut != nil, "unpack destination %v for type %v is nil", i, prettyType{t})
		assertf(ut.Kind() == reflect.Ptr, "unpack destination %v for type %v is %v, not a pointer", i, prettyType{t}, ut)

		ft := t.Field(i).Type
		assertf(ft.AssignableTo(ut.Elem()), "%v from %v is not assignable to unpack destination %v of type %v", prettyType{ft}, prettyType{t}, i, prettyType{ut.Elem()})

		reflect.ValueOf(u).Elem().Set(v.Field(i))
	}
}

type prettyValue struct {
	Value
}

func (v prettyValue) String() string {
	t := reflect.TypeOf(v.Value)

	if isTuple(t) {
		var fields []string
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, fmt.Sprintf("%#v", reflect.ValueOf(v.Value).Field(i).Interface()))
		}

		return "(" + strings.Join(fields, ", ") + ")"
	} else {
		return fmt.Sprintf("%#v", v.Value)
	}
}

type prettyType struct {
	reflect.Type
}

func (t prettyType) String() string {
	if isTuple(t) {
		var fields []string
		for i := 0; i < t.NumField(); i++ {
			fields = append(fields, t.Field(i).Type.String())
		}

		return "(" + strings.Join(fields, ", ") + ")"
	} else {
		return t.Type.String()
	}
}
