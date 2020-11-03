// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"testing"

	. "pgregory.net/rapid"
)

func TestURL(t *testing.T) {
	pathEscapeRegex := regexp.MustCompile(`^[0-9A-Fa-f]{2}`)

	Check(t, func(t *T) {
		u := URL().Draw(t, "url").(url.URL)

		// should be parseable
		if _, err := url.Parse(u.String()); err != nil {
			t.Fatalf("URL returned unparseable url %s: %v", u.String(), err)
		}

		// only valid characters in path
		for i, ch := range u.Path {
			if !(('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || strings.ContainsRune("$-_.+!*'(),%/@=&:~", ch)) {
				t.Fatalf("URL returned invalid url %s: invalid character %s at %d", u.String(), string(ch), i)
			}
		}

		// assert proper path escapes
		for _, co := range strings.Split(u.Path, "%")[1:] {
			if ok := pathEscapeRegex.MatchString(co); !ok {
				t.Fatalf("URL returned invalid url %s: invalid escape %s", u.String(), co)
			}
		}
	})
}

func TestDomainOf(t *testing.T) {
	t.Parallel()

	genFuncs := []func(int, int) *Generator{
		func(i, j int) *Generator { return DomainOf(i, j) },
	}

	for i, gf := range genFuncs {
		t.Run(strconv.Itoa(i), MakeCheck(func(t *T) {
			maxLength := IntRange(4, 255).Draw(t, "maxLength").(int)
			maxElementLength := IntRange(1, 63).Draw(t, "maxElementLength").(int)

			d := gf(maxLength, maxElementLength).Draw(t, "d").(string)
			if got, want := len(d), maxLength; got > want {
				t.Errorf("got domain of length %d with maxLenght of %d", got, want)
			}

			elements := strings.Split(d, ".")

			// ignore the tld
			for i, elem := range elements[:len(elements)-1] {
				if got, want := len(elem), maxElementLength; got > want {
					t.Errorf("got domain element %d of length %d with maxElementLength %d", i, got, want)
				}
			}
		}))
	}
}
