// Copyright 2020 Walter Scheper <walter.scheper@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"unicode"
)

const (
	domainMaxLength        = 255
	domainMaxElementLength = 63
)

var (
	domainType = reflect.TypeOf("")
	urlType    = reflect.TypeOf(url.URL{})
)

type domainNameGen struct{}

func (*domainNameGen) String() string {
	return "Domain()"
}

var tldGenerator = SampledFrom(tlds)

func (g *domainNameGen) value(t *T) string {
	domain := tldGenerator.
		Filter(func(s string) bool { return len(s)+2 <= domainMaxLength }).
		Draw(t, "domain")

	expr := fmt.Sprintf(`[a-zA-Z]([a-zA-Z0-9\-]{0,%d}[a-zA-Z0-9])?`, domainMaxElementLength-2)
	elements := newRepeat(1, 126, 1, g.String())
	for elements.more(t.s) {
		subDomain := StringMatching(expr).Draw(t, "subdomain")
		if len(domain)+len(subDomain) >= domainMaxLength {
			break
		}
		domain = subDomain + "." + domain
	}

	return domain
}

// Domain generates an RFC 1035 compliant domain name.
func Domain() *Generator[string] {
	return newGenerator[string](&domainNameGen{})
}

type urlGenerator struct {
	schemes []string
}

func (g *urlGenerator) String() string {
	return "URL()"
}

var printableGen = StringOf(RuneFrom(nil, unicode.PrintRanges...))

func (g *urlGenerator) value(t *T) url.URL {
	scheme := SampledFrom(g.schemes).Draw(t, "scheme")
	var domain string
	switch SampledFrom([]int{0, 1, 2}).Draw(t, "g") {
	case 2:
		domain = Domain().Draw(t, "domain")
	case 1:
		domain = IPv6().Draw(t, "domain").String()
		domain = "[" + domain + "]"
	case 0:
		domain = IPv4().Draw(t, "domain").String()
	}
	port := IntRange(0, 2^16-1).Draw(t, "port")
	path_ := SliceOf(printableGen).Draw(t, "path")
	query := SliceOf(printableGen).Draw(t, "query")
	fragment := printableGen.Draw(t, "fragment")

	// join domain and port
	if port > 0 {
		domain += fmt.Sprintf(":%d", port)
	}

	// URL escape path
	for i := range path_ {
		path_[i] = url.PathEscape(path_[i])
	}

	// url escape query strings
	for i := range query {
		query[i] = url.QueryEscape(query[i])
	}

	return url.URL{
		Host:     domain,
		Path:     strings.Join(path_, "/"),
		Scheme:   scheme,
		RawQuery: strings.Join(query, "&"),
		Fragment: fragment,
	}
}

// URL generates RFC 3986 compliant http/https URLs.
func URL() *Generator[url.URL] {
	return urlOf([]string{"http", "https"})
}

func urlOf(schemes []string) *Generator[url.URL] {
	return newGenerator[url.URL](&urlGenerator{
		schemes: schemes,
	})
}
