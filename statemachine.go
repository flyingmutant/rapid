// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
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
	actionLabel      = "action"
	validActionTries = 100 // hack, but probably good enough for now

	initMethodPrefix  = "Init"
	checkMethodName   = "Check"
	cleanupMethodName = "Cleanup"

	noValidActionsMsg = "can't find a valid action"
)

// StateMachine is a convenience function for defining "state machine" tests,
// to be run by Check or MakeCheck.
//
// State machine test is a pattern for testing stateful systems that looks
// like this:
//
//   s := new(STATE_MACHINE_DEFINITION_TYPE)
//   RUN_RANDOM_INITIALIZER_ACTION(s)
//   defer CLEANUP(s)
//   CHECK_INVARIANTS(s)
//   for {
//       RUN_RANDOM_ACTION(s)
//       CHECK_INVARIANTS(s)
//   }
//
// StateMachine synthesizes such test from the type of its argument,
// which must be a pointer to an arbitrary state machine definition type,
// whose public methods are treated as follows:
//
// - Init(t *rapid.T), InitAnySuffixHere(t *rapid.T), if present,
// are used as "initializer" actions; exactly one is ran at the beginning
// of each test case;
//
// - Cleanup(), if present, is called at the end of each test case;
//
// - Check(t *rapid.T), if present, is ran after every action;
//
// - All other public methods should have a form ActionName(t *rapid.T)
// and are used as possible actions. At least one action has to be specified.
func StateMachine(i interface{}) func(*T) {
	typ := reflect.TypeOf(i)

	return func(t *T) {
		t.Helper()

		repeat := newRepeat(0, *steps, maxInt)

		sm := newStateMachine(typ)
		sm.init(t)
		defer sm.cleanup()

		sm.checkInvariants(t)
		for repeat.more(t.s, typ.String()) {
			ok := sm.executeAction(t)
			if ok {
				sm.checkInvariants(t)
			} else {
				repeat.reject()
			}
		}
	}
}

type stateMachine struct {
	inits      map[string]func(*T)
	actions    map[string]func(*T)
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
		inits      = map[string]func(*T){}
		actions    = map[string]func(*T){}
		initKeys   []string
		actionKeys []string
		check      func(*T)
		cleanup    func()
	)

	for i := 0; i < n; i++ {
		name := typ.Method(i).Name
		m, ok := v.Method(i).Interface().(func(*T))
		if ok {
			if strings.HasPrefix(name, initMethodPrefix) {
				inits[name] = m
				initKeys = append(initKeys, name)
			} else if name == checkMethodName {
				check = m
			} else {
				actions[name] = m
				actionKeys = append(actionKeys, name)
			}
		} else {
			assertf(name == cleanupMethodName, "unexpected state machine method %v", name)
			m, ok := v.Method(i).Interface().(func())
			assertf(ok, "method %v should have type func(), not %v", cleanupMethodName, v.Method(i).Type())
			cleanup = m
		}
	}

	assertf(len(actions) > 0, "state machine of type %v has no actions specified", typ)

	sort.Strings(initKeys)
	sort.Strings(actionKeys)

	sm := &stateMachine{
		inits:      inits,
		actions:    actions,
		actionKeys: SampledFrom(actionKeys),
		check:      check,
		cleanup_:   cleanup,
	}
	if len(initKeys) > 0 {
		sm.initKeys = SampledFrom(initKeys)
	}

	return sm
}

func (sm *stateMachine) init(t *T) {
	if sm.initKeys != nil {
		t.Helper()
		sm.inits[sm.initKeys.Draw(t, "initializer").(string)](t)
	}
}

func (sm *stateMachine) cleanup() {
	if sm.cleanup_ != nil {
		sm.cleanup_()
	}
}

func (sm *stateMachine) executeAction(t *T) bool {
	t.Helper()

	for n := 0; n < validActionTries; n++ {
		i := t.s.beginGroup(actionLabel, false)
		action := sm.actions[sm.actionKeys.Draw(t, "action").(string)]
		invalid, skipped := runAction(t, action)
		t.s.endGroup(i, false)

		if skipped {
			continue
		} else {
			return !invalid
		}
	}

	panic(stopTest(noValidActionsMsg))
}

func (sm *stateMachine) checkInvariants(t *T) {
	if sm.check != nil {
		t.Helper()
		sm.check(t)
	}
}

func runAction(t *T, action func(*T)) (invalid bool, skipped bool) {
	defer func(draws int) {
		if r := recover(); r != nil {
			if _, ok := r.(invalidData); ok {
				invalid = true
				skipped = t.draws == draws
			} else {
				panic(r)
			}
		}
	}(t.draws)

	action(t)

	return false, false
}
