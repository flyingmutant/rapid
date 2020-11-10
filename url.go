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

func (*domainNameGen) type_() reflect.Type {
	return domainType
}

var tldGenerator = SampledFrom(tlds)

func (g *domainNameGen) value(t *T) value {
	domain := tldGenerator.
		Filter(func(s string) bool { return len(s)+2 <= domainMaxLength }).
		Map(func(s string) string {
			var n string
			for _, ch := range s {
				n += string(SampledFrom([]rune{unicode.ToLower(ch), unicode.ToUpper(ch)}).Draw(t, "").(rune))
			}

			return n
		}).
		Draw(t, "domain").(string)

	expr := fmt.Sprintf(`[a-zA-Z]([a-zA-Z0-9\-]{0,%d}[a-zA-Z0-9])?`, domainMaxElementLength-2)
	elements := newRepeat(1, 126, 1)
	for elements.more(t.s, g.String()) {
		subDomain := StringMatching(expr).Draw(t, "subdomain").(string)
		if len(domain)+len(subDomain) >= domainMaxLength {
			break
		}
		domain = subDomain + "." + domain
	}

	return domain
}

// Domain generates an RFC 1035 compliant domain name.
func Domain() *Generator {
	return newGenerator(&domainNameGen{})
}

type urlGenerator struct {
	schemes []string
}

func (g *urlGenerator) String() string {
	return "URL()"
}

func (g *urlGenerator) type_() reflect.Type {
	return urlType
}

var printableGen = StringOf(RuneFrom(nil, unicode.PrintRanges...))

func (g *urlGenerator) value(t *T) value {
	scheme := SampledFrom(g.schemes).Draw(t, "scheme").(string)
	domain := Domain().Draw(t, "domain").(string)
	port := IntRange(0, 2^16-1).
		Map(func(i int) string {
			if i == 0 {
				return ""
			}
			return fmt.Sprintf(":%d", i)
		}).
		Draw(t, "port").(string)
	path_ := SliceOf(printableGen).Draw(t, "path").([]string)
	query := SliceOf(printableGen.Map(url.QueryEscape)).Draw(t, "query").([]string)
	fragment := printableGen.Draw(t, "fragment").(string)

	return url.URL{
		Host:   domain + port,
		Path:     strings.Join(path_, "/"),
		Scheme: scheme,
		RawQuery: strings.Join(query, "&"),
		Fragment: fragment,
	}
}

// URL generates RFC 3986 compliant http/https URLs.
func URL() *Generator {
	return urlOf([]string{"http", "https"})
}

func urlOf(schemes []string) *Generator {
	return newGenerator(&urlGenerator{
		schemes: schemes,
	})
}
