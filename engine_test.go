// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"strings"
	"testing"
)

func brokenGen(t *T) int { panic("this generator is not working") }

type brokenMachine struct{}

func (m *brokenMachine) DoNothing(*T) { panic("this state machine is not working") }
func (m *brokenMachine) Check(t *T)   {}

func TestPanicTraceback(t *testing.T) {
	t.Parallel()

	testData := []struct {
		name   string
		suffix string
		fail   func(*T) *testError
	}{
		{
			"impossible filter",
			"pgregory.net/rapid.find",
			func(t *T) *testError {
				g := Bool().Filter(func(bool) bool { return false })
				_, err := recoverValue(g, t)
				return err
			},
		},
		{
			"broken custom generator",
			"pgregory.net/rapid.brokenGen",
			func(t *T) *testError {
				g := Custom(brokenGen)
				_, err := recoverValue(g, t)
				return err
			},
		},
		{
			"broken state machine",
			"pgregory.net/rapid.(*brokenMachine).DoNothing",
			func(t *T) *testError {
				return checkOnce(t, Run(&brokenMachine{}))
			},
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			s := createRandomBitStream(t)
			nt := newT(t, s, false)

			err := td.fail(nt)
			if err == nil {
				t.Fatalf("test case did not fail")
			}

			lines := strings.Split(err.traceback, "\n")
			if !strings.HasSuffix(lines[0], td.suffix) {
				t.Errorf("bad traceback:\n%v", err.traceback)
			}
		})
	}
}

func BenchmarkCheckOverhead(b *testing.B) {
	g := Uint()
	f := func(t *T) {
		g.Draw(t, "")
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		checkTB(b, f)
	}
}
