// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"fmt"
	"regexp"
	"testing"

	. "pgregory.net/rapid"
)

type tEvent struct {
	t      *testing.T
	output []string
}

func NewTEvent(t *testing.T) *tEvent {
	return &tEvent{
		t:      t,
		output: make([]string, 0),
	}
}

// implement rapid.TB trivially except for Logf
func (t *tEvent) Helper()                                   { t.t.Helper() }
func (t *tEvent) Name() string                              { return t.t.Name() }
func (t *tEvent) Error(args ...interface{})                 { t.t.Error(args...) }
func (t *tEvent) Errorf(format string, args ...interface{}) { t.t.Errorf(format, args...) }
func (t *tEvent) Fail()                                     { t.t.Fail() }
func (t *tEvent) FailNow()                                  { t.t.FailNow() }
func (t *tEvent) Failed() bool                              { return t.t.Failed() }
func (t *tEvent) Fatal(args ...interface{})                 { t.t.Fatal(args...) }
func (t *tEvent) Fatalf(format string, args ...interface{}) { t.t.Fatalf(format, args...) }
func (t *tEvent) SkipNow()                                  { t.t.SkipNow() }
func (t *tEvent) Skipped() bool                             { return t.t.Skipped() }
func (t *tEvent) Log(args ...interface{})                   { t.t.Log(args...) }
func (t *tEvent) Logf(format string, data ...interface{}) {
	t.t.Logf(format, data...)
	t.output = append(t.output, fmt.Sprintf(format, data...))
}

func TestEventEmitter(t *testing.T) {
	te := NewTEvent(t)
	Event(te, "x")
	Event(te, "y")

	PrintStats(te)
	checkMatch(t, fmt.Sprintf("Statistics.*%s", t.Name()), te.output[0])
	checkMatch(t, "of 2 ", te.output[1])
	checkMatch(t, "x: 1 \\(50.0+ %", te.output[3])
	checkMatch(t, "y: 1 \\(50.0+ %", te.output[4])
}

func checkMatch(t *testing.T, pattern, str string) {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		t.Fatalf("Regex compile failed")
	}
	if !matched {
		t.Fatalf("Pattern <%s> does not match in <%s>", pattern, str)
	}
}
