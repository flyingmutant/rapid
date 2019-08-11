// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"strconv"
	"testing"
	"unicode"
	"unicode/utf8"

	. "github.com/flyingmutant/rapid"
)

func TestStringsExamples(t *testing.T) {
	g := StringsN(10, -1, -1)

	for i := 0; i < 100; i++ {
		s, _, _ := g.Example()
		t.Log(len(s.(string)), s)
	}
}

func TestRegexpExamples(t *testing.T) {
	g := StringsMatching("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	for i := 0; i < 100; i++ {
		s, _, _ := g.Example()
		t.Log(s)
	}
}

func TestStringsOfRunesAreUTF8(t *testing.T) {
	t.Parallel()

	gens := []*Generator{
		Strings(),
		StringsN(2, 10, -1),
		StringsOf(Runes()),
		StringsOfN(Runes(), 2, 10, -1),
		StringsOf(RunesFrom(nil, unicode.Cyrillic)),
		StringsOf(RunesFrom([]rune{'a', 'b', 'c'})),
	}

	for _, g := range gens {
		t.Run(g.String(), MakeCheck(func(t *T) {
			s := g.Draw(t, "s").(string)
			if !utf8.ValidString(s) {
				t.Fatalf("invalid UTF-8 string: %q", s)
			}
		}))
	}
}

func TestStringRuneCountLimits(t *testing.T) {
	t.Parallel()

	genFuncs := []func(i, j int) *Generator{
		func(i, j int) *Generator { return StringsN(i, j, -1) },
		func(i, j int) *Generator { return StringsOfN(Runes(), i, j, -1) },
	}

	for i, gf := range genFuncs {
		t.Run(strconv.Itoa(i), MakeCheck(func(t *T) {
			minRunes := IntsRange(0, 256).Draw(t, "minRunes").(int)
			maxRunes := IntsMin(0).Draw(t, "maxRunes").(int)
			if minRunes > maxRunes {
				t.Skip("minRunes > maxRunes")
			}

			s := gf(minRunes, maxRunes).Draw(t, "s").(string)
			n := utf8.RuneCountInString(s)
			if n < minRunes {
				t.Fatalf("got string with %v runes with lower limit %v", n, minRunes)
			}
			if n > maxRunes {
				t.Fatalf("got string with %v runes with upper limit %v", n, maxRunes)
			}
		}))
	}
}

func TestStringsNMaxLen(t *testing.T) {
	t.Parallel()

	genFuncs := []func(int) *Generator{
		func(i int) *Generator { return StringsN(-1, -1, i) },
		func(i int) *Generator { return StringsOfN(Runes(), -1, -1, i) },
		func(i int) *Generator { return StringsOfNBytes(-1, i) },
	}

	for i, gf := range genFuncs {
		t.Run(strconv.Itoa(i), MakeCheck(func(t *T) {
			maxLen := IntsMin(0).Draw(t, "maxLen").(int)
			s := gf(maxLen).Draw(t, "s").(string)
			if len(s) > maxLen {
				t.Fatalf("got string of length %v with maxLen %v", len(s), maxLen)
			}
		}))
	}
}
