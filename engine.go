// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"bytes"
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
	tracebackStop = "pgregory.net/rapid.checkOnce"
	runtimePrefix = "runtime."
)

var (
	flags cmdline

	emptyStructType  = reflect.TypeOf(struct{}{})
	emptyStructValue = reflect.ValueOf(struct{}{})

	tracebackBlacklist = map[string]bool{
		"pgregory.net/rapid.(*customGen).maybeValue.func1": true,
		"pgregory.net/rapid.runAction.func1":               true,
	}
)

type cmdline struct {
	checks     int
	steps      int
	failfile   string
	nofailfile bool
	seed       uint64
	log        bool
	verbose    bool
	debug      bool
	debugvis   bool
	shrinkTime time.Duration
}

func init() {
	flag.IntVar(&flags.checks, "rapid.checks", 100, "rapid: number of checks to perform")
	flag.IntVar(&flags.steps, "rapid.steps", 100, "rapid: number of state machine steps to perform")
	flag.StringVar(&flags.failfile, "rapid.failfile", "", "rapid: fail file to use to reproduce test failure")
	flag.BoolVar(&flags.nofailfile, "rapid.nofailfile", false, "rapid: do not write fail files on test failures")
	flag.Uint64Var(&flags.seed, "rapid.seed", 0, "rapid: PRNG seed to start with (0 to use a random one)")
	flag.BoolVar(&flags.log, "rapid.log", false, "rapid: eager verbose output to stdout (to aid with unrecoverable test failures)")
	flag.BoolVar(&flags.verbose, "rapid.v", false, "rapid: verbose output")
	flag.BoolVar(&flags.debug, "rapid.debug", false, "rapid: debugging output")
	flag.BoolVar(&flags.debugvis, "rapid.debugvis", false, "rapid: debugging visualization")
	flag.DurationVar(&flags.shrinkTime, "rapid.shrinktime", 30*time.Second, "rapid: maximum time to spend on test case minimization")
}

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
	valid, invalid, seed, buf, err1, err2 := doCheck(tb, flags.failfile, flags.checks, baseSeed(), prop)
	dt := time.Since(start)

	if err1 == nil && err2 == nil {
		if valid == flags.checks {
			tb.Logf("[rapid] OK, passed %v tests (%v)", valid, dt)
		} else {
			tb.Errorf("[rapid] only generated %v valid tests from %v total (%v)", valid, valid+invalid, dt)
		}
	} else {
		repr := fmt.Sprintf("-rapid.seed=%d", seed)
		if flags.failfile != "" && seed == 0 {
			repr = fmt.Sprintf("-rapid.failfile=%q", flags.failfile)
		} else if !flags.nofailfile {
			_, failfile := failFileName(tb.Name())
			out := captureTestOutput(tb, prop, buf)
			err := saveFailFile(failfile, rapidVersion, out, seed, buf)
			if err == nil {
				repr = fmt.Sprintf("-rapid.failfile=%q (or -rapid.seed=%d)", failfile, seed)
			} else {
				tb.Logf("[rapid] %v", err)
			}
		}

		name := regexp.QuoteMeta(tb.Name())
		if traceback(err1) == traceback(err2) {
			if err2.isStopTest() {
				tb.Errorf("[rapid] failed after %v tests: %v\nTo reproduce, specify -run=%q %v\nFailed test output:", valid, err2, name, repr)
			} else {
				tb.Errorf("[rapid] panic after %v tests: %v\nTo reproduce, specify -run=%q %v\nTraceback:\n%vFailed test output:", valid, err2, name, repr, traceback(err2))
			}
		} else {
			tb.Errorf("[rapid] flaky test, can not reproduce a failure\nTo try to reproduce, specify -run=%q %v\nTraceback (%v):\n%vOriginal traceback (%v):\n%vFailed test output:", name, repr, err2, traceback(err2), err1, traceback(err1))
		}

		_ = checkOnce(newT(tb, newBufBitStream(buf, false), true, nil), prop) // output using (*testing.T).Log for proper line numbers
	}

	if tb.Failed() {
		tb.FailNow() // do not try to run any checks after the first failed one
	}
}

func doCheck(tb tb, failfile string, checks int, seed uint64, prop func(*T)) (int, int, uint64, []uint64, *testError, *testError) {
	tb.Helper()

	assertf(!tb.Failed(), "check function called with *testing.T which has already failed")

	if failfile != "" {
		buf, err1, err2 := checkFailFile(tb, failfile, prop)
		if err1 != nil || err2 != nil {
			return 0, 0, 0, buf, err1, err2
		}
	}

	seed, valid, invalid, err1 := findBug(tb, checks, seed, prop)
	if err1 == nil {
		return valid, invalid, 0, nil, nil, nil
	}

	s := newRandomBitStream(seed, true)
	t := newT(tb, s, flags.verbose, nil)
	t.Logf("[rapid] trying to reproduce the failure")
	err2 := checkOnce(t, prop)
	if !sameError(err1, err2) {
		return valid, invalid, seed, s.data, err1, err2
	}

	t.Logf("[rapid] trying to minimize the failing test case")
	buf, err3 := shrink(tb, s.recordedBits, err2, prop)

	return valid, invalid, seed, buf, err2, err3
}

func checkFailFile(tb tb, failfile string, prop func(*T)) ([]uint64, *testError, *testError) {
	tb.Helper()

	version, _, buf, err := loadFailFile(failfile)
	if err != nil {
		tb.Logf("[rapid] ignoring fail file: %v", err)
		return nil, nil, nil
	}
	if version != rapidVersion {
		tb.Logf("[rapid] ignoring fail file: version %q differs from rapid version %q", version, rapidVersion)
		return nil, nil, nil
	}

	s1 := newBufBitStream(buf, false)
	t1 := newT(tb, s1, flags.verbose, nil)
	err1 := checkOnce(t1, prop)
	if err1 == nil {
		return nil, nil, nil
	}
	if err1.isInvalidData() {
		tb.Logf("[rapid] fail file %q is no longer valid", failfile)
		return nil, nil, nil
	}

	s2 := newBufBitStream(buf, false)
	t2 := newT(tb, s2, flags.verbose, nil)
	t2.Logf("[rapid] trying to reproduce the failure")
	err2 := checkOnce(t2, prop)

	return buf, err1, err2
}

func findBug(tb tb, checks int, seed uint64, prop func(*T)) (uint64, int, int, *testError) {
	tb.Helper()

	var (
		r       = newRandomBitStream(0, false)
		t       = newT(tb, r, flags.verbose, nil)
		valid   = 0
		invalid = 0
	)

	for valid < checks && invalid < checks*invalidChecksMult {
		seed += uint64(valid) + uint64(invalid)
		r.init(seed)
		var start time.Time
		if t.shouldLog() {
			t.Logf("[rapid] test #%v start (seed %v)", valid+invalid+1, seed)
			start = time.Now()
		}

		err := checkOnce(t, prop)
		if err == nil {
			if t.shouldLog() {
				t.Logf("[rapid] test #%v OK (%v)", valid+invalid+1, time.Since(start))
			}
			valid++
		} else if err.isInvalidData() {
			if t.shouldLog() {
				t.Logf("[rapid] test #%v invalid (%v)", valid+invalid+1, time.Since(start))
			}
			invalid++
		} else {
			if t.shouldLog() {
				t.Logf("[rapid] test #%v failed: %v", valid+invalid+1, err)
			}
			return seed, valid, invalid, err
		}
	}

	return 0, valid, invalid, nil
}

func checkOnce(t *T, prop func(*T)) (err *testError) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	defer func() { err = panicToError(recover(), 3) }()

	prop(t)
	t.failOnError()

	return nil
}

func captureTestOutput(tb tb, prop func(*T), buf []uint64) []byte {
	var b bytes.Buffer
	l := log.New(&b, fmt.Sprintf("%s ", tb.Name()), log.Ldate|log.Ltime) // TODO: enable log.Lmsgprefix once all supported versions of Go have it
	_ = checkOnce(newT(tb, newBufBitStream(buf, false), false, l), prop)
	return b.Bytes()
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
	f, more, skipSpecial := runtime.Frame{}, true, true
	for more && !strings.HasSuffix(f.Function, tracebackStop) {
		f, more = frames.Next()

		if skipSpecial && (tracebackBlacklist[f.Function] || strings.HasPrefix(f.Function, runtimePrefix)) {
			continue
		}
		skipSpecial = false

		_, err := fmt.Fprintf(b, "    %s:%d in %s\n", f.File, f.Line, f.Function)
		assert(err == nil)
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

// TB is a common interface between *testing.T, *testing.B and *T.
type TB interface {
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

// tb is a private copy of TB, made to avoid T having public fields
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

// T is similar to testing.T, but with extra bookkeeping for property-based tests.
//
// For tests to be reproducible, they should generally run in a single goroutine.
// If concurrency is unavoidable, methods on *T, such as Helper and Errorf, are safe for concurrent calls,
// but Draw from a given *T is not.
type T struct {
	tb       // unnamed to force re-export of (*T).Helper()
	tbLog    bool
	rawLog   *log.Logger
	s        bitStream
	draws    int
	refDraws []value
	mu       sync.RWMutex
	failed   stopTest
}

func newT(tb tb, s bitStream, tbLog bool, rawLog *log.Logger, refDraws ...value) *T {
	t := &T{
		tb:       tb,
		tbLog:    tbLog,
		rawLog:   rawLog,
		s:        s,
		refDraws: refDraws,
	}

	if rawLog == nil && flags.log {
		testName := "rapid test"
		if tb != nil {
			testName = tb.Name()
		}

		t.rawLog = log.New(os.Stdout, fmt.Sprintf("[%v] ", testName), 0)
	}

	return t
}

func (t *T) draw(g *Generator, label string) value {
	v := g.value(t)

	if len(t.refDraws) > 0 {
		ref := t.refDraws[t.draws]
		if !reflect.DeepEqual(v, ref) {
			t.tb.Fatalf("draw %v differs: %#v vs expected %#v", t.draws, v, ref)
		}
	}

	if t.tbLog || t.rawLog != nil {
		if label == "" {
			label = fmt.Sprintf("#%v", t.draws)
		}

		if t.tbLog && t.tb != nil {
			t.tb.Helper()
		}
		t.Logf("[rapid] draw %v: %#v", label, v)
	}

	t.draws++

	return v
}

func (t *T) shouldLog() bool {
	return t.rawLog != nil || (t.tbLog && t.tb != nil)
}

func (t *T) Logf(format string, args ...interface{}) {
	if t.rawLog != nil {
		t.rawLog.Printf(format, args...)
	} else if t.tbLog && t.tb != nil {
		t.tb.Helper()
		t.tb.Logf(format, args...)
	}
}

func (t *T) Log(args ...interface{}) {
	if t.rawLog != nil {
		t.rawLog.Print(args...)
	} else if t.tbLog && t.tb != nil {
		t.tb.Helper()
		t.tb.Log(args...)
	}
}

// Skipf is equivalent to Logf followed by SkipNow.
func (t *T) Skipf(format string, args ...interface{}) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.skip(fmt.Sprintf(format, args...))
}

// Skip is equivalent to Log followed by SkipNow.
func (t *T) Skip(args ...interface{}) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	t.Log(args...)
	t.skip(fmt.Sprint(args...))
}

// SkipNow marks the current test case as invalid (except state machine
// tests, where it marks current action as non-applicable instead).
// If too many test cases are skipped, rapid will mark the test as failing
// due to inability to generate enough valid test cases.
//
// Prefer Filter to SkipNow, and prefer generators that always produce
// valid test cases to Filter.
func (t *T) SkipNow() {
	t.skip("(*T).SkipNow() called")
}

// Errorf is equivalent to Logf followed by Fail.
func (t *T) Errorf(format string, args ...interface{}) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.fail(false, fmt.Sprintf(format, args...))
}

// Error is equivalent to Log followed by Fail.
func (t *T) Error(args ...interface{}) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	t.Log(args...)
	t.fail(false, fmt.Sprint(args...))
}

// Fatalf is equivalent to Logf followed by FailNow.
func (t *T) Fatalf(format string, args ...interface{}) {
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	t.Logf(format, args...)
	t.fail(true, fmt.Sprintf(format, args...))
}

// Fatal is equivalent to Log followed by FailNow.
func (t *T) Fatal(args ...interface{}) {
	if t.tbLog && t.tb != nil {
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

func (t *T) failOnError() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.failed != "" {
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
