// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	small             = 5
	invalidChecksMult = 10
	exampleMaxTries   = 1000

	tracebackLen  = 32
	tracebackStop = "github.com/flyingmutant/rapid.checkOnce"
)

var (
	checks     = flag.Int("rapid.checks", 100, "rapid: number of checks to perform")
	steps      = flag.Int("rapid.steps", 100, "rapid: number of state machine steps to perform")
	startSeed  = flag.Uint64("rapid.seed", 0, "rapid: PRNG seed to start with (0 to use a random one)")
	rapidLog   = flag.Bool("rapid.log", false, "rapid: eager verbose output to stdout (to aid with unrecoverable test failures)")
	verbose    = flag.Bool("rapid.v", false, "rapid: verbose output")
	debug      = flag.Bool("rapid.debug", false, "rapid: debugging output")
	debugvis   = flag.Bool("rapid.debugvis", false, "rapid: debugging visualization")
	shrinkTime = flag.Duration("rapid.shrinktime", 30*time.Second, "rapid: maximum time to spend on test case minimization")

	errCantGenDueToFilter = errors.New("generation failed due to Filter() or Assume() conditions being too strong")

	emptyStructType  = reflect.TypeOf(struct{}{})
	emptyStructValue = reflect.ValueOf(struct{}{})
)

func assert(ok bool) {
	if !ok {
		panic("assertion failed")
	}
}

func assertf(ok bool, format string, args ...interface{}) {
	if !ok {
		panic(fmt.Sprintf(format, args...))
	}
}

func assertValidRange(min int, max int) {
	assertf(max < 0 || min <= max, fmt.Sprintf("invalid range [%d, %d]", min, max))
}

// Check fails the current test if rapid can find a test case which falsifies prop.
//
// Property is falsified in case of a panic or a call to
// (*T).Fatalf, (*T).Fatal, (*T).Errorf, (*T).Error, (*T).FailNow or (*T).Fail.
func Check(t *testing.T, prop func(*T)) {
	t.Helper()
	checkTB(t, prop)
}

// MakeCheck is a convenience function for defining subtests suitable for
// (*testing.T).Run. It allows you to write this:
//
//   t.Run("subtest name", rapid.MakeCheck(func(t *rapid.T) {
//       // test code
//   }))
//
// instead of this:
//
//   t.Run("subtest name", func(t *testing.T) {
//       rapid.Check(t, func(t *rapid.T) {
//           // test code
//       })
//   })
//
func MakeCheck(prop func(*T)) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		checkTB(t, prop)
	}
}

func checkTB(tb tb, prop func(*T)) {
	tb.Helper()

	start := time.Now()
	valid, invalid, seed, buf, err1, err2 := doCheck(tb, prop)
	dt := time.Since(start)

	if err1 == nil && err2 == nil {
		if valid == *checks {
			tb.Logf("[rapid] OK, passed %v tests (%v)", valid, dt)
		} else {
			tb.Errorf("[rapid] only generated %v valid tests from %v total (%v)", valid, valid+invalid, dt)
		}
	} else {
		name := regexp.QuoteMeta(tb.Name())
		if traceback(err1) == traceback(err2) {
			if err2.isStopTest() {
				tb.Errorf("[rapid] failed after %v tests: %v\nTo reproduce, specify -run=%q -rapid.seed=%v\nFailed test output:", valid, err2, name, seed)
			} else {
				tb.Errorf("[rapid] panic after %v tests: %v\nTo reproduce, specify -run=%q -rapid.seed=%v\nTraceback:\n%vFailed test output:", valid, err2, name, seed, traceback(err2))
			}
		} else {
			tb.Errorf("[rapid] flaky test, can not reproduce a failure\nTo try to reproduce, specify -run=%q -rapid.seed=%v\nTraceback (%v):\n%vOriginal traceback (%v):\n%vFailed test output:", name, seed, err2, traceback(err2), err1, traceback(err1))
		}

		_ = checkOnce(newT(tb, newBufBitStream(buf, false), true), prop)
	}

	if tb.Failed() {
		tb.FailNow() // do not try to run any checks after the first failed one
	}
}

func doCheck(tb tb, prop func(*T)) (int, int, uint64, []uint64, *testError, *testError) {
	tb.Helper()

	assertf(!tb.Failed(), "check function called with *testing.T which has already failed")

	seed, valid, invalid, err1 := findBug(tb, baseSeed(), prop)
	if err1 == nil {
		return valid, invalid, 0, nil, nil, nil
	}

	s := newRandomBitStream(seed, true)
	t := newT(tb, s, *verbose)
	t.Logf("[rapid] trying to reproduce the failure")
	err2 := checkOnce(t, prop)
	if !sameError(err1, err2) {
		return valid, invalid, seed, s.data, err1, err2
	}

	t.Logf("[rapid] trying to minimize the failing test case")
	buf, err3 := shrink(tb, s.recordedBits, err2, prop)

	return valid, invalid, seed, buf, err2, err3
}

func findBug(tb tb, seed uint64, prop func(*T)) (uint64, int, int, *testError) {
	tb.Helper()

	valid := 0
	invalid := 0
	for valid < *checks && invalid < *checks*invalidChecksMult {
		start := time.Now()
		seed += uint64(valid) + uint64(invalid)
		t := newT(tb, newRandomBitStream(seed, false), *verbose)
		t.Logf("[rapid] test #%v start (seed %v)", valid+invalid+1, seed)
		err := checkOnce(t, prop)
		dt := time.Since(start)
		if err == nil {
			t.Logf("[rapid] test #%v OK (%v)", valid+invalid+1, dt)
			valid++
		} else if err.isInvalidData() {
			t.Logf("[rapid] test #%v invalid (%v)", valid+invalid+1, dt)
			invalid++
		} else {
			t.Logf("[rapid] test #%v failed: %v", valid+invalid+1, err)
			return seed, valid, invalid, err
		}
	}

	return 0, valid, invalid, nil
}

func checkOnce(t *T, prop func(*T)) (err *testError) {
	t.Helper()

	defer func() { err = panicToError(recover(), 3) }()

	prop(t)

	if t.Failed() {
		panic(t.failed)
	}

	return nil
}

type invalidData string
type stopTest string

type testError struct {
	data      interface{}
	traceback string
}

func panicToError(p interface{}, skip int) *testError {
	if p == nil {
		return nil
	}

	callers := make([]uintptr, tracebackLen)
	callers = callers[:runtime.Callers(skip, callers)]
	frames := runtime.CallersFrames(callers)

	b := &strings.Builder{}
	f, more, skipRuntime := runtime.Frame{}, true, true
	for more && f.Function != tracebackStop {
		f, more = frames.Next()

		isRuntime := strings.HasPrefix(f.Function, "runtime.")
		if !isRuntime {
			skipRuntime = false
		}
		if !isRuntime || !skipRuntime {
			_, err := fmt.Fprintf(b, "    %s:%d in %s\n", f.File, f.Line, f.Function)
			assert(err == nil)
		}
	}

	return &testError{
		data:      p,
		traceback: b.String(),
	}
}

func (err *testError) Error() string {
	if msg, ok := err.data.(stopTest); ok {
		return string(msg)
	}

	if msg, ok := err.data.(invalidData); ok {
		return fmt.Sprintf("invalid data: %s", string(msg))
	}

	return fmt.Sprintf("%v", err.data)
}

func (err *testError) isInvalidData() bool {
	_, ok := err.data.(invalidData)
	return ok
}

func (err *testError) isStopTest() bool {
	_, ok := err.data.(stopTest)
	return ok
}

func sameError(err1 *testError, err2 *testError) bool {
	return errorString(err1) == errorString(err2) && traceback(err1) == traceback(err2)
}

func errorString(err *testError) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func traceback(err *testError) string {
	if err == nil {
		return "    <no error>\n"
	}

	return err.traceback
}

type tb interface {
	Helper()
	Name() string
	Logf(format string, args ...interface{})
	Log(args ...interface{})
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	FailNow()
	Fail()
	Failed() bool
}

type T struct {
	tb       // unnamed to force re-export of (*T).Helper()
	log      bool
	rapidLog *log.Logger
	s        bitStream
	draws    int
	refDraws []value
	mu       sync.RWMutex
	failed   stopTest
}

func newT(tb tb, s bitStream, log_ bool, refDraws ...value) *T {
	t := &T{
		tb:       tb,
		log:      log_,
		s:        s,
		refDraws: refDraws,
	}

	if *rapidLog {
		testName := "rapid test"
		if tb != nil {
			testName = tb.Name()
		}

		t.rapidLog = log.New(os.Stdout, fmt.Sprintf("[%v] ", testName), 0)
	}

	return t
}

func (t *T) draw(g *Generator, label string) value {
	v := g.value(t)

	if len(t.refDraws) > 0 {
		ref := t.refDraws[t.draws]
		if !reflect.DeepEqual(v, ref) {
			t.tb.Fatalf("draw %v differs: %v vs expected %v", t.draws, prettyValue{v}, prettyValue{ref})
		}
	}

	if t.log || t.rapidLog != nil {
		if label == "" {
			label = fmt.Sprintf("#%v", t.draws)
		}

		if t.tb != nil {
			t.tb.Helper()
		}
		t.Logf("[rapid] draw %v: %v", label, prettyValue{v})
	}

	t.draws++

	return v
}

func (t *T) Logf(format string, args ...interface{}) {
	if t.rapidLog != nil {
		t.rapidLog.Printf(format, args...)
	} else if t.log && t.tb != nil {
		t.tb.Helper()
		t.tb.Logf(format, args...)
	}
}

func (t *T) Log(args ...interface{}) {
	if t.rapidLog != nil {
		t.rapidLog.Print(args...)
	} else if t.log && t.tb != nil {
		t.tb.Helper()
		t.tb.Log(args...)
	}
}

func (t *T) Skipf(format string, args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.skip(fmt.Sprintf(format, args...))
}

func (t *T) Skip(args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Log(args...)
	t.skip(fmt.Sprint(args...))
}

// SkipNow marks the current test case as invalid.
// If too many test cases are skipped, rapid will mark the test as failing
// due to inability to generate enough valid test cases.
//
// Prefer Filter to SkipNow, and prefer generators that always produce
// valid test cases to Filter.
func (t *T) SkipNow() {
	t.skip("(*T).SkipNow() called")
}

func (t *T) Errorf(format string, args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.fail(false, fmt.Sprintf(format, args...))
}

func (t *T) Error(args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Log(args...)
	t.fail(false, fmt.Sprint(args...))
}

func (t *T) Fatalf(format string, args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.fail(true, fmt.Sprintf(format, args...))
}

func (t *T) Fatal(args ...interface{}) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.Log(args...)
	t.fail(true, fmt.Sprint(args...))
}

func (t *T) FailNow() {
	t.fail(true, "(*T).FailNow() called")
}

func (t *T) Fail() {
	t.fail(false, "(*T).Fail() called")
}

func (t *T) Failed() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.failed != ""
}

func (t *T) skip(msg string) {
	panic(invalidData(msg))
}

func (t *T) fail(now bool, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.failed = stopTest(msg)
	if now {
		panic(t.failed)
	}
}

func assertCallable(fn reflect.Type, t reflect.Type, name string) {
	assertf(fn.Kind() == reflect.Func, "%v should be a function, not %v", name, fn.Kind())
	assertf(fn.NumIn() == 1, "%v should have 1 parameter, not %v", name, fn.NumIn())
	assertf(fn.NumOut() == 1, "%v should have 1 output parameter, not %v", name, fn.NumOut())
	assertf(t.AssignableTo(fn.In(0)), "parameter #0 (%v) of %v should be assignable from %v", fn.In(0), name, t)
}

func call(fn reflect.Value, arg reflect.Value) value {
	r := fn.Call([]reflect.Value{arg})

	if len(r) == 0 {
		return nil
	} else {
		assert(len(r) == 1)
		return r[0].Interface()
	}
}

type prettyValue struct {
	value
}

func (v prettyValue) String() string {
	return fmt.Sprintf("%#v", v.value)
}
