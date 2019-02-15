// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"regexp/syntax"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

const runesGenString = "Runes()"

var (
	stringType    = reflect.TypeOf("")
	byteSliceType = reflect.TypeOf([]byte(nil))

	defaultRunes = []rune{
		'?',
		'~', '!', '@', '#', '$', '%', '^', '&', '*', '_', '-', '+', '=',
		'.', ',', ':', ';',
		' ', '\t', '\r', '\n',
		'/', '\\', '|',
		'(', '[', '{', '<',
		'\'', '"', '`',
		'\x00', '\x0B', '\x1B', '\x7F',
		'\uFEFF', '\uFFFD', '\u202E',
		'Ⱥ',
	}

	// unicode.Categories without surrogates (which are not allowed in UTF-8), ordered by taste
	defaultTables = []*unicode.RangeTable{
		unicode.Lu, // Letter, uppercase        (1781)
		unicode.Ll, // Letter, lowercase        (2145)
		unicode.Lt, // Letter, titlecase          (31)
		unicode.Lm, // Letter, modifier          (250)
		unicode.Lo, // Letter, other          (121212)
		unicode.Nd, // Number, decimal digit     (610)
		unicode.Nl, // Number, letter            (236)
		unicode.No, // Number, other             (807)
		unicode.P,  // Punctuation               (788)
		unicode.Sm, // Symbol, math              (948)
		unicode.Sc, // Symbol, currency           (57)
		unicode.Sk, // Symbol, modifier          (121)
		unicode.So, // Symbol, other            (5984)
		unicode.Mn, // Mark, nonspacing         (1805)
		unicode.Me, // Mark, enclosing            (13)
		unicode.Mc, // Mark, spacing combining   (415)
		unicode.Z,  // Separator                  (19)
		unicode.Cc, // Other, control             (65)
		unicode.Cf, // Other, format             (152)
		unicode.Co, // Other, private use     (137468)
	}

	expandedTablesMu = sync.RWMutex{}
	expandedTables   = map[string][]rune{}
)

func Runes() *Generator {
	return runesFrom(true, defaultRunes, defaultTables...)
}

func RunesFrom(runes []rune, tables ...*unicode.RangeTable) *Generator {
	return runesFrom(false, runes, tables...)
}

func runesFrom(default_ bool, runes []rune, tables ...*unicode.RangeTable) *Generator {
	if len(tables) == 0 {
		assertf(len(runes) > 0, "at least one rune should be specified")
	}
	if len(runes) == 0 {
		assertf(len(tables) > 0, "at least one *unicode.RangeTable should be specified")
	}

	var weights []int
	if len(runes) > 0 {
		weights = append(weights, len(tables))
	}
	for range tables {
		weights = append(weights, 1)
	}

	tables_ := make([][]rune, len(tables))
	for i := range tables {
		tables_[i] = expandRangeTable(tables[i], rangeTableName(tables[i]))
		assertf(len(tables_[i]) > 0, "empty *unicode.RangeTable %v", i)
	}

	return newGenerator(&runeGen{
		die:      newLoadedDie(weights),
		runes:    runes,
		tables:   tables_,
		default_: default_,
	})
}

type runeGen struct {
	die      *loadedDie
	runes    []rune
	tables   [][]rune
	default_ bool
}

func (g *runeGen) String() string {
	if g.default_ {
		return runesGenString
	} else {
		return fmt.Sprintf("Runes(%v runes, %v tables)", len(g.runes), len(g.tables))
	}
}

func (g *runeGen) type_() reflect.Type {
	return int32Type
}

func (g *runeGen) value(s bitStream) Value {
	n := g.die.roll(s)

	runes := g.runes
	if len(g.runes) == 0 {
		runes = g.tables[n]
	} else if n > 0 {
		runes = g.tables[n-1]
	}

	return runes[genIndex(s, len(runes), true)]
}

func Strings() *Generator {
	return StringsOf(Runes())
}

func StringsN(minRunes int, maxRunes int, maxLen int) *Generator {
	return StringsOfN(Runes(), minRunes, maxRunes, maxLen)
}

func StringsOf(rune_ *Generator) *Generator {
	return StringsOfN(rune_, -1, -1, -1)
}

func StringsOfN(rune_ *Generator, minRunes int, maxRunes int, maxLen int) *Generator {
	assertValidRange(minRunes, maxRunes)
	assertf(rune_.type_() == int32Type, "rune generator should generate %v, not %v", int32Type, prettyType{rune_.type_()})
	assertf(maxLen < 0 || maxLen >= maxRunes, "maximum length (%v) should not be less than maximum number of runes (%v)", maxLen, maxRunes)

	return newGenerator(&stringGen{
		runeGen:  rune_,
		minElems: minRunes,
		maxElems: maxRunes,
		maxLen:   maxLen,
	})
}

func StringsOfBytes() *Generator {
	return StringsOfNBytes(-1, -1)
}

func StringsOfNBytes(minLen int, maxLen int) *Generator {
	assertValidRange(minLen, maxLen)

	return newGenerator(&stringGen{
		byteGen:  Bytes(),
		minElems: minLen,
		maxElems: maxLen,
		maxLen:   -1,
	})
}

type stringGen struct {
	byteGen  *Generator
	runeGen  *Generator
	minElems int
	maxElems int
	maxLen   int
}

func (g *stringGen) String() string {
	if g.runeGen != nil && g.runeGen.String() == runesGenString {
		if g.minElems < 0 && g.maxElems < 0 && g.maxLen < 0 {
			return "Strings()"
		} else {
			return fmt.Sprintf("StringsN(minRunes=%v, maxRunes=%v, maxLen=%v)", g.minElems, g.maxElems, g.maxLen)
		}
	} else if g.runeGen != nil {
		if g.minElems < 0 && g.maxElems < 0 && g.maxLen < 0 {
			return fmt.Sprintf("StringsOf(%v)", g.runeGen)
		} else {
			return fmt.Sprintf("StringsOfN(%v, minRunes=%v, maxRunes=%v, maxLen=%v)", g.runeGen, g.minElems, g.maxElems, g.maxLen)
		}
	} else {
		if g.minElems < 0 && g.maxElems < 0 && g.maxLen < 0 {
			return "StringsOfBytes()"
		} else {
			return fmt.Sprintf("StringsOfNBytes(minLen=%v, maxLen=%v)", g.minElems, g.maxElems)
		}
	}
}

func (g *stringGen) type_() reflect.Type {
	return stringType
}

func (g *stringGen) value(s bitStream) Value {
	repeat := newRepeat(g.minElems, g.maxElems, -1)

	var b strings.Builder
	b.Grow(repeat.avg())

	if g.runeGen != nil {
		maxLen := g.maxLen
		if maxLen < 0 {
			maxLen = maxInt
		}

		for repeat.more(s, g.runeGen.String()) {
			r := g.runeGen.value(s).(rune)
			n := utf8.RuneLen(r)

			if n < 0 || b.Len()+n > maxLen {
				repeat.reject()
			} else {
				b.WriteRune(r)
			}
		}
	} else {
		for repeat.more(s, g.byteGen.String()) {
			b.WriteByte(g.byteGen.value(s).(byte))
		}
	}

	return b.String()
}

func StringsMatching(expr string) *Generator {
	return matching(expr, true)
}

func SlicesOfBytesMatching(expr string) *Generator {
	return matching(expr, false)
}

func matching(expr string, str bool) *Generator {
	syn, err := syntax.Parse(expr, syntax.Perl)
	assertf(err == nil, "failed to parse regexp %q: %v", expr, err)

	re, err := regexp.Compile(expr)
	assertf(err == nil, "failed to compile regexp %q: %v", expr, err)

	return newGenerator(&regexpGen{
		str:         str,
		expr:        expr,
		syn:         syn,
		re:          re,
		any:         Runes(),
		anyNoNL:     Runes().Filter(func(r rune) bool { return r != '\n' }),
		subNames:    map[*syntax.Regexp]string{},
		charClasses: map[*syntax.Regexp]*Generator{},
	})
}

type runeWriter interface {
	WriteRune(r rune) (int, error)
}

type regexpGen struct {
	str         bool
	expr        string
	syn         *syntax.Regexp
	re          *regexp.Regexp
	any         *Generator
	anyNoNL     *Generator
	subNames    map[*syntax.Regexp]string
	charClasses map[*syntax.Regexp]*Generator
}

func (g *regexpGen) String() string {
	if g.str {
		return fmt.Sprintf("StringsMatching(%q)", g.expr)
	} else {
		return fmt.Sprintf("SlicesOfBytesMatching(%q)", g.expr)
	}
}

func (g *regexpGen) type_() reflect.Type {
	if g.str {
		return stringType
	} else {
		return byteSliceType
	}
}

func (g *regexpGen) value(s bitStream) Value {
	var (
		gen    func(bitStream) Value
		filter func(Value) bool
	)

	if g.str {
		gen = func(s bitStream) Value {
			b := &strings.Builder{}
			g.build(b, g.syn, s)
			return b.String()
		}

		filter = func(v Value) bool {
			return g.re.MatchString(v.(string))
		}
	} else {
		gen = func(s bitStream) Value {
			b := &bytes.Buffer{}
			g.build(b, g.syn, s)
			return b.Bytes()
		}

		filter = func(v Value) bool {
			return g.re.Match(v.([]byte))
		}
	}

	return satisfy(filter, gen, s, small, "")
}

func (g *regexpGen) build(w runeWriter, re *syntax.Regexp, s bitStream) {
	i := s.beginGroup(re.Op.String(), false)

	switch re.Op {
	case syntax.OpNoMatch:
		panic(invalidData("no possible regexp match"))
	case syntax.OpEmptyMatch:
		s.drawBits(0)
	case syntax.OpLiteral:
		s.drawBits(0)
		for _, r := range re.Rune {
			w.WriteRune(maybeFoldCase(s, r, re.Flags))
		}
	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		sub := g.any
		switch re.Op {
		case syntax.OpCharClass:
			sub = g.loadGen(re)
		case syntax.OpAnyCharNotNL:
			sub = g.anyNoNL
		}
		r := sub.value(s).(rune)
		w.WriteRune(maybeFoldCase(s, r, re.Flags))
	case syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText,
		syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		s.drawBits(0) // do nothing and hope that Assume() is enough
	case syntax.OpCapture:
		g.build(w, re.Sub[0], s)
	case syntax.OpStar, syntax.OpPlus, syntax.OpQuest, syntax.OpRepeat:
		min, max := re.Min, re.Max
		switch re.Op {
		case syntax.OpStar:
			min, max = 0, -1
		case syntax.OpPlus:
			min, max = 1, -1
		case syntax.OpQuest:
			min, max = 0, 1
		}
		repeat := newRepeat(min, max, -1)
		for repeat.more(s, g.loadName(re.Sub[0])) {
			g.build(w, re.Sub[0], s)
		}
	case syntax.OpConcat:
		for _, sub := range re.Sub {
			g.build(w, sub, s)
		}
	case syntax.OpAlternate:
		ix := genIndex(s, len(re.Sub), true)
		g.build(w, re.Sub[ix], s)
	default:
		assertf(false, "invalid regexp op %v", re.Op)
	}

	s.endGroup(i, false)
}

func (g *regexpGen) loadName(re *syntax.Regexp) string {
	name := g.subNames[re]
	if name == "" {
		name = re.String()
		g.subNames[re] = name
	}
	return name
}

func (g *regexpGen) loadGen(re *syntax.Regexp) *Generator {
	sub := g.charClasses[re]
	if sub == nil {
		sub = charClassGen(re)
		g.charClasses[re] = sub
	}
	return sub
}

func maybeFoldCase(s bitStream, r rune, flags syntax.Flags) rune {
	n := 0
	if flags&syntax.FoldCase != 0 {
		n = int(genUintN(s, 4, false))
	}

	for i := 0; i < n; i++ {
		r = unicode.SimpleFold(r)
	}

	return r
}

func rangeTableName(t *unicode.RangeTable) string {
	maps := map[string]map[string]*unicode.RangeTable{
		"cat":        unicode.Categories,
		"foldcat":    unicode.FoldCategory,
		"foldscript": unicode.FoldScript,
		"prop":       unicode.Properties,
		"script":     unicode.Scripts,
	}

	for c, m := range maps {
		for k, v := range m {
			if v == t {
				return fmt.Sprintf("%s/%s", c, k)
			}
		}
	}

	return fmt.Sprintf("%p", t)
}

func expandRangeTable(t *unicode.RangeTable, name string) []rune {
	expandedTablesMu.RLock()
	ret, ok := expandedTables[name]
	expandedTablesMu.RUnlock()
	if ok {
		return ret
	}

	expandedTablesMu.Lock()
	defer expandedTablesMu.Unlock()
	ret, ok = expandedTables[name]
	if ok {
		return ret
	}

	for _, r := range t.R16 {
		for i := r.Lo; i <= r.Hi; i += r.Stride {
			ret = append(ret, rune(i))
		}
	}
	for _, r := range t.R32 {
		for i := r.Lo; i <= r.Hi; i += r.Stride {
			ret = append(ret, rune(i))
		}
	}
	expandedTables[name] = ret

	return ret
}

func charClassGen(re *syntax.Regexp) *Generator {
	assert(re.Op == syntax.OpCharClass)

	t := &unicode.RangeTable{}
	for i := 0; i < len(re.Rune); i += 2 {
		t.R32 = append(t.R32, unicode.Range32{
			Lo:     uint32(re.Rune[i]),
			Hi:     uint32(re.Rune[i+1]),
			Stride: 1,
		})
	}

	return newGenerator(&runeGen{
		die:    newLoadedDie([]int{1}),
		tables: [][]rune{expandRangeTable(t, re.String())},
	})
}
