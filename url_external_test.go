// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"net/url"
	"strings"
	"testing"

	. "pgregory.net/rapid"
)

func TestURL(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T) {
		u := URL().Draw(t, "url").(url.URL)

		// should be parseable
		if _, err := url.Parse(u.String()); err != nil {
			t.Fatalf("URL returned unparseable url %s: %v", u.String(), err)
		}
	})
}

func TestDomain(t *testing.T) {
	t.Parallel()

	Check(t, func(t *T) {
		d := Domain().Draw(t, "d").(string)
		if got, want := len(d), 255; got > want {
			t.Errorf("got domain of length %d with maxLenght of %d", got, want)
		}

		elements := strings.Split(d, ".")

		// ignore the tld
		for i, elem := range elements[:len(elements)-1] {
			if got, want := len(elem), 63; got > want {
				t.Errorf("got domain element %d of length %d with maxElementLength %d", i, got, want)
			}
		}
	})
}
