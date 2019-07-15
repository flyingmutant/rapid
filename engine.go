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

	tPtrType         = reflect.TypeOf((*T)(nil))
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

// Assume marks the current test case as invalid if cond is false.
// If assumption is too hard to satisfy, rapid will mark the test as failing
// due to inability to generate enough valid test cases.
//
// Prefer Filter to Assume, and prefer generators that always produce
// valid test cases to Filter.
func Assume(cond bool) {
	if !cond {
		panic(invalidData("failed to satisfy assumption"))
	}
}

// Bind is a convenience function for writing state machine transition rules.
//
// Given prop with signature func(*T, Arg1,  ..., ArgN) and N generators
// producing values of types Arg1, ...,  ArgN, Bind returns a function
// of type func(*T) which calls prop with the values drawn from the
// generators as arguments 1...N.
func Bind(prop interface{}, args ...*Generator) func(*T) {
	if len(args) == 0 {
		fn, ok := prop.(func(*T))
		assertf(ok, "prop should have type func(*T), not %v", reflect.TypeOf(prop))
		return fn
	}

	var (
		args_ = Tuple(args...)
		pv    = reflect.ValueOf(prop)
		pt    = reflect.TypeOf(prop)
	)

	assertCallable(pt, args_.type_(), "prop", 1)
	assertf(pt.In(0) == tPtrType, "prop should have first parameter of type %v, not %v", tPtrType, pt.In(0))
	assertf(pt.NumOut() == 0, "prop should have no output parameters (got %v)", pt.NumOut())

	return func(t *T) {
		t.Helper()

		v := reflect.ValueOf(args_.Draw(t, "args"))

		n := v.NumField()
		in := make([]reflect.Value, n+1)
		in[0] = reflect.ValueOf(t)

		for i := 0; i < n; i++ {
			in[i+1] = v.Field(i)
		}

		pv.Call(in)
	}
}

// BindIf is a convenience function for writing conditional state machine
// transition rules.
//
// BindIf behaves exactly like Bind, except that it returns nil if
// the precondition is false.
func BindIf(precondition bool, prop interface{}, args ...*Generator) func(*T) {
	if !precondition {
		return nil
	}

	return Bind(prop, args...)
}

// Check fails the current test if rapid can find a test case which falsifies
// the property created by Bind(prop, args...). Property is falsified in case
// of a panic or a call to (*T).Fatalf, (*T).Fatal, (*T).Errorf, (*T).Error,
// (*T).FailNow or (*T).Fail.
func Check(t *testing.T, prop interface{}, args ...*Generator) {
	t.Helper()
	checkTB(t, Bind(prop, args...))
}

// MakeCheck is a convenience function for defining subtests suitable for
// (*testing.T).Run.
func MakeCheck(prop interface{}, args ...*Generator) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()
		checkTB(t, Bind(prop, args...))
	}
}

func checkTB(tb limitedTB, prop func(*T)) {
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

func doCheck(tb limitedTB, prop func(*T)) (int, int, uint64, []uint64, *testError, *testError) {
	tb.Helper()

	assertf(!tb.Failed(), "check function called with *testing.T which has already failed")

	seed, valid, invalid, err1 := findBug(tb, prngSeed(), prop)
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

func findBug(tb limitedTB, seed uint64, prop func(*T)) (uint64, int, int, *testError) {
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

type limitedTB interface {
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
	limitedTB // unnamed to force re-export of (*T).Helper()
	log       bool
	rapidLog  *log.Logger
	src       *bitStreamSource
	draws     int
	refDraws  []Value
	mu        sync.RWMutex
	failed    stopTest
}

func newT(tb limitedTB, s bitStream, log_ bool, refDraws ...Value) *T {
	t := &T{
		limitedTB: tb,
		log:       log_,
		src:       &bitStreamSource{s},
		refDraws:  refDraws,
	}

	if *rapidLog {
		t.rapidLog = log.New(os.Stdout, fmt.Sprintf("[%v] ", tb.Name()), 0)
	}

	return t
}

func (t *T) draw(g *Generator, label string, unpack ...interface{}) Value {
	v := t.src.draw(g, label, unpack...)

	if len(t.refDraws) > 0 {
		ref := t.refDraws[t.draws]
		if !reflect.DeepEqual(v, ref) {
			t.limitedTB.Fatalf("draw %v differs: %v vs expected %v", t.draws, prettyValue{v}, prettyValue{ref})
		}
	}

	if t.log || t.rapidLog != nil {
		if label == "" {
			label = fmt.Sprintf("#%v", t.draws)
		}

		t.Helper()
		t.Logf("[rapid] draw %v: %v", label, prettyValue{v})
	}

	t.draws++

	return v
}

func (t *T) Logf(format string, args ...interface{}) {
	if t.rapidLog != nil {
		t.rapidLog.Printf(format, args...)
	} else if t.log {
		t.Helper()
		t.limitedTB.Logf(format, args...)
	}
}

func (t *T) Log(args ...interface{}) {
	if t.rapidLog != nil {
		t.rapidLog.Print(args...)
	} else if t.log {
		t.Helper()
		t.limitedTB.Log(args...)
	}
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.Helper()
	t.Logf(format, args...)
	t.fail(false, fmt.Sprintf(format, args...))
}

func (t *T) Error(args ...interface{}) {
	t.Helper()
	t.Log(args...)
	t.fail(false, fmt.Sprint(args...))
}

func (t *T) Fatalf(format string, args ...interface{}) {
	t.Helper()
	t.Logf(format, args...)
	t.fail(true, fmt.Sprintf(format, args...))
}

func (t *T) Fatal(args ...interface{}) {
	t.Helper()
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

func (t *T) fail(now bool, msg string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.failed = stopTest(msg)
	if now {
		panic(t.failed)
	}
}
