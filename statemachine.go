// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"reflect"
	"sort"
	"strings"
)

const (
	validActionTries = 100 // hack, but probably good enough for now

	initMethodPrefix  = "Init"
	checkMethodName   = "Check"
	cleanupMethodName = "Cleanup"

	noValidActionsMsg = "can't find a valid action"
)

// StateMachine synthesizes a property to be checked with Check or MakeCheck
// from the type of its argument, which must be a pointer to a state machine
// definition type.
func StateMachine(i interface{}) func(*T) {
	typ := reflect.TypeOf(i)

	return func(t *T) {
		t.Helper()

		repeat := newRepeat(0, *steps, maxInt)

		sm := newStateMachine(typ)
		sm.init(t)
		defer sm.cleanup()

		sm.checkInvariants(t)
		for repeat.more(t.data.s, typ.String()) {
			sm.selectAction(t)(t)
			sm.checkInvariants(t)
		}
	}
}

type stateMachine struct {
	inits      map[string]func() func(*T)
	actions    map[string]func() func(*T)
	initKeys   *Generator
	actionKeys *Generator
	check      func(*T)
	cleanup_   func()
}

func newStateMachine(typ reflect.Type) *stateMachine {
	assertf(typ.Kind() == reflect.Ptr, "state machine type should be a pointer, not %v", typ.Kind())

	var (
		v          = reflect.New(typ.Elem())
		n          = typ.NumMethod()
		inits      = map[string]func() func(*T){}
		actions    = map[string]func() func(*T){}
		initKeys   []string
		actionKeys []string
		check      func(*T)
		cleanup_   func()
	)

	for i := 0; i < n; i++ {
		m, ok := v.Method(i).Interface().(func() func(*T))
		if ok {
			name := typ.Method(i).Name

			if strings.HasPrefix(name, initMethodPrefix) {
				inits[name] = m
				initKeys = append(initKeys, name)
			} else {
				actions[name] = m
				actionKeys = append(actionKeys, name)
			}
		}
	}

	if checkM := v.MethodByName(checkMethodName); checkM.IsValid() {
		check, _ = checkM.Interface().(func(*T))
		assertf(check != nil, "method %v should have type func(*T), not %v", checkMethodName, checkM.Type())
	}

	if cleanupM := v.MethodByName(cleanupMethodName); cleanupM.IsValid() {
		cleanup_, _ = cleanupM.Interface().(func())
		assertf(cleanup_ != nil, "method %v should have type func(), not %v", cleanupMethodName, cleanupM.Type())
	}

	assertf(len(actions) > 0, "state machine of type %v has no actions specified", typ)

	sort.Strings(initKeys)
	sort.Strings(actionKeys)

	sm := &stateMachine{
		inits:      inits,
		actions:    actions,
		actionKeys: filter(SampledFrom(actionKeys), func(key string) bool { return actions[key]() != nil }, validActionTries, noValidActionsMsg),
		check:      check,
		cleanup_:   cleanup_,
	}
	if len(initKeys) > 0 {
		sm.initKeys = SampledFrom(initKeys)
	}

	return sm
}

func (sm *stateMachine) init(t *T) {
	if sm.initKeys != nil {
		t.Helper()
		sm.inits[sm.initKeys.Draw(t, "initializer").(string)]()(t)
	}
}

func (sm *stateMachine) cleanup() {
	if sm.cleanup_ != nil {
		sm.cleanup_()
	}
}

func (sm *stateMachine) selectAction(t *T) func(*T) {
	t.Helper()

	return sm.actions[sm.actionKeys.Draw(t, "action").(string)]()
}

func (sm *stateMachine) checkInvariants(t *T) {
	if sm.check != nil {
		t.Helper()
		sm.check(t)
	}
}
