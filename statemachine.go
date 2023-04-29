// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"math"
	"sort"
	"testing"
)

const (
	actionLabel      = "action"
	validActionTries = 100 // hack, but probably good enough for now

	noValidActionsMsg = "can't find a valid action"
)

// Run executes a random sequence of actions (often called a "state machine" test).
// actions[""], if set, is executed before/after every other action invocation
// and should only contain invariant checking code.
func (t *T) Run(actions map[string]func(*T)) {
	t.Helper()
	if len(actions) == 0 {
		return
	}

	check := func(*T) {}
	actionKeys := make([]string, 0, len(actions))
	for key, action := range actions {
		if key != "" {
			actionKeys = append(actionKeys, key)
		} else {
			check = action
		}
	}
	sort.Strings(actionKeys)

	steps := flags.steps
	if testing.Short() {
		steps /= 5
	}

	repeat := newRepeat(0, steps, math.MaxInt, "Run")
	sm := stateMachine{
		check:      check,
		actionKeys: SampledFrom(actionKeys),
		actions:    actions,
	}

	sm.check(t)
	t.failOnError()
	for repeat.more(t.s) {
		ok := sm.executeAction(t)
		if ok {
			sm.check(t)
			t.failOnError()
		} else {
			repeat.reject()
		}
	}
}

type stateMachine struct {
	check      func(*T)
	actionKeys *Generator[string]
	actions    map[string]func(*T)
}

func (sm *stateMachine) executeAction(t *T) bool {
	t.Helper()

	for n := 0; n < validActionTries; n++ {
		i := t.s.beginGroup(actionLabel, false)
		action := sm.actions[sm.actionKeys.Draw(t, "action")]
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
	t.failOnError()

	return false, false
}
