package rapid

import (
	"fmt"
	"net/url"
	"path"
	"reflect"
	"strings"
	"unicode"
)

type domainNameGen struct {
	maxLength        int
	maxElementLength int
}

func (g *domainNameGen) String() string {
	return fmt.Sprintf("Domain(maxLength=%v, mmaxElementLength%v)", g.maxLength, g.maxElementLength)
}

func (g *domainNameGen) type_() reflect.Type {
	return stringType
}

func (g *domainNameGen) value(t *T) value {
	domain := SampledFrom(tlds).
		Filter(func(s string) bool { return len(s)+2 <= g.maxLength }).
		Map(func(s string) string {
			var n string
			for _, ch := range s {
				n += string(SampledFrom([]rune{unicode.ToUpper(ch), unicode.ToLower(ch)}).Draw(t, "").(rune))
			}

			return n
		}).Draw(t, "domain").(string)

	var b strings.Builder
	elements := newRepeat(1, 126, 1)
	b.Grow(elements.avg())

	var expr string
	switch g.maxElementLength {
	case 1:
		expr = `[a-zA-Z]`
	case 2:
		expr = `[a-zA-Z][a-zA-Z0-9]?`
	default:
		expr = fmt.Sprintf(`[a-zA-Z]([a-zA-Z0-9\-]{0,%d}[a-zA-Z0-9])?`, g.maxElementLength-2)
	}
	for elements.more(t.s, g.String()) {
		subDomain := StringMatching(expr).Draw(t, "subdomain").(string)
		if len(domain)+len(subDomain) >= g.maxLength {
			break
		}
		domain = subDomain + "." + domain
	}

	return domain
}

// Domain generates an RFC 1035 compliant domain name.
func Domain() *Generator {
	return DomainOf(255, 63)
}

// DomainOf generates an RFC 1035 compliant domain name,
// with a maximum overall length of maxLength
// and a maximum number of elements of maxElements.
func DomainOf(maxLength, maxElementLength int) *Generator {
	assertf(4 <= maxLength, "maximum length (%v) should not be less than 4, to generate a two character domain and a one character subdomain", maxLength)
	assertf(maxLength <= 255, "maximum length (%v) should not be greater than 255 to comply with RFC 1035", maxLength)
	assertf(1 <= maxElementLength, "maximum element length (%v) should not be less than 1 to comply with RFC 1035", maxElementLength)
	assertf(maxElementLength <= 63, "maximum element length (%v) should not be greater than 63 to comply with RFC 1035", maxElementLength)

	return newGenerator(&domainNameGen{
		maxElementLength: maxElementLength,
		maxLength:        maxLength,
	})
}

type urlGenerator struct {
	schemes []string
}

func (g *urlGenerator) String() string {
	return fmt.Sprintf("URLGenerator(schemes=%v)", g.schemes)
}

func (g *urlGenerator) type_() reflect.Type {
	return reflect.TypeOf(url.URL{})
}

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
	path_ := path.Join(
		SliceOf(
			StringOf(RuneFrom(nil, unicode.PrintRanges...)).Map(url.PathEscape),
		).Draw(t, "path").([]string)...)

	return url.URL{
		Host:   domain + port,
		Path:   path_,
		Scheme: scheme,
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
