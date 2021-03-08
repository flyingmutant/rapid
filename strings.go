// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
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

var (
	stringType    = reflect.TypeOf("")
	byteSliceType = reflect.TypeOf([]byte(nil))

	defaultRunes = []rune{
		'A', 'a', '?',
		'~', '!', '@', '#', '$', '%', '^', '&', '*', '_', '-', '+', '=',
		'.', ',', ':', ';',
		' ', '\t', '\r', '\n',
		'/', '\\', '|',
		'(', '[', '{', '<',
		'\'', '"', '`',
		'\x00', '\x0B', '\x1B', '\x7F', // NUL, VT, ESC, DEL
		'\uFEFF', '\uFFFD', '\u202E', // BOM, replacement character, RTL override
		'Ⱥ', // In UTF-8, Ⱥ increases in length from 2 to 3 bytes when lowercased
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

	expandedTables  = sync.Map{} // *unicode.RangeTable / regexp name -> []rune
	compiledRegexps = sync.Map{} // regexp -> compiledRegexp
	regexpNames     = sync.Map{} // *regexp.Regexp -> string
	charClassGens   = sync.Map{} // regexp name -> *Generator

	anyRuneGen     = Rune()
	anyRuneGenNoNL = Rune().Filter(func(r rune) bool { return r != '\n' })
)

type compiledRegexp struct {
	syn *syntax.Regexp
	re  *regexp.Regexp
}

func Rune() *Generator {
	return runesFrom(true, defaultRunes, defaultTables...)
}

func RuneFrom(runes []rune, tables ...*unicode.RangeTable) *Generator {
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
		tables_[i] = expandRangeTable(tables[i], tables[i])
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
		return "Rune()"
	} else {
		return fmt.Sprintf("Rune(%v runes, %v tables)", len(g.runes), len(g.tables))
	}
}

func (g *runeGen) type_() reflect.Type {
	return int32Type
}

func (g *runeGen) value(t *T) value {
	n := g.die.roll(t.s)

	runes := g.runes
	if len(g.runes) == 0 {
		runes = g.tables[n]
	} else if n > 0 {
		runes = g.tables[n-1]
	}

	return runes[genIndex(t.s, len(runes), true)]
}

func String() *Generator {
	return StringOf(anyRuneGen)
}

func StringN(minRunes int, maxRunes int, maxLen int) *Generator {
	return StringOfN(anyRuneGen, minRunes, maxRunes, maxLen)
}

func StringOf(elem *Generator) *Generator {
	return StringOfN(elem, -1, -1, -1)
}

func StringOfN(elem *Generator, minElems int, maxElems int, maxLen int) *Generator {
	assertValidRange(minElems, maxElems)
	assertf(elem.type_() == int32Type || elem.type_() == uint8Type, "element generator should generate runes or bytes, not %v", elem.type_())
	assertf(maxLen < 0 || maxLen >= maxElems, "maximum length (%v) should not be less than maximum number of elements (%v)", maxLen, maxElems)

	return newGenerator(&stringGen{
		elem:     elem,
		minElems: minElems,
		maxElems: maxElems,
		maxLen:   maxLen,
	})
}

type stringGen struct {
	elem     *Generator
	minElems int
	maxElems int
	maxLen   int
}

func (g *stringGen) String() string {
	if g.elem == anyRuneGen {
		if g.minElems < 0 && g.maxElems < 0 && g.maxLen < 0 {
			return "String()"
		} else {
			return fmt.Sprintf("StringN(minRunes=%v, maxRunes=%v, maxLen=%v)", g.minElems, g.maxElems, g.maxLen)
		}
	} else {
		if g.minElems < 0 && g.maxElems < 0 && g.maxLen < 0 {
			return fmt.Sprintf("StringOf(%v)", g.elem)
		} else {
			return fmt.Sprintf("StringOfN(%v, minElems=%v, maxElems=%v, maxLen=%v)", g.elem, g.minElems, g.maxElems, g.maxLen)
		}
	}
}

func (g *stringGen) type_() reflect.Type {
	return stringType
}

func (g *stringGen) value(t *T) value {
	repeat := newRepeat(g.minElems, g.maxElems, -1)

	var b strings.Builder
	b.Grow(repeat.avg())

	if g.elem.type_() == int32Type {
		maxLen := g.maxLen
		if maxLen < 0 {
			maxLen = maxInt
		}

		for repeat.more(t.s, g.elem.String()) {
			r := g.elem.value(t).(rune)
			n := utf8.RuneLen(r)

			if n < 0 || b.Len()+n > maxLen {
				repeat.reject()
			} else {
				b.WriteRune(r)
			}
		}
	} else {
		for repeat.more(t.s, g.elem.String()) {
			b.WriteByte(g.elem.value(t).(byte))
		}
	}

	return b.String()
}

func StringMatching(expr string) *Generator {
	return matching(expr, true)
}

func SliceOfBytesMatching(expr string) *Generator {
	return matching(expr, false)
}

func matching(expr string, str bool) *Generator {
	compiled, err := compileRegexp(expr)
	assertf(err == nil, "%v", err)

	return newGenerator(&regexpGen{
		str:  str,
		expr: expr,
		syn:  compiled.syn,
		re:   compiled.re,
	})
}

type runeWriter interface {
	WriteRune(r rune) (int, error)
}

type regexpGen struct {
	str  bool
	expr string
	syn  *syntax.Regexp
	re   *regexp.Regexp
}

func (g *regexpGen) String() string {
	if g.str {
		return fmt.Sprintf("StringMatching(%q)", g.expr)
	} else {
		return fmt.Sprintf("SliceOfBytesMatching(%q)", g.expr)
	}
}

func (g *regexpGen) type_() reflect.Type {
	if g.str {
		return stringType
	} else {
		return byteSliceType
	}
}

func (g *regexpGen) maybeString(t *T) value {
	b := &strings.Builder{}
	g.build(b, g.syn, t)
	v := b.String()

	if g.re.MatchString(v) {
		return v
	} else {
		return nil
	}
}

func (g *regexpGen) maybeSlice(t *T) value {
	b := &bytes.Buffer{}
	g.build(b, g.syn, t)
	v := b.Bytes()

	if g.re.Match(v) {
		return v
	} else {
		return nil
	}
}

func (g *regexpGen) value(t *T) value {
	if g.str {
		return find(g.maybeString, t, small)
	} else {
		return find(g.maybeSlice, t, small)
	}
}

func (g *regexpGen) build(w runeWriter, re *syntax.Regexp, t *T) {
	i := t.s.beginGroup(re.Op.String(), false)

	switch re.Op {
	case syntax.OpNoMatch:
		panic(invalidData("no possible regexp match"))
	case syntax.OpEmptyMatch:
		t.s.drawBits(0)
	case syntax.OpLiteral:
		t.s.drawBits(0)
		for _, r := range re.Rune {
			_, _ = w.WriteRune(maybeFoldCase(t.s, r, re.Flags))
		}
	case syntax.OpCharClass, syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		sub := anyRuneGen
		switch re.Op {
		case syntax.OpCharClass:
			sub = charClassGen(re)
		case syntax.OpAnyCharNotNL:
			sub = anyRuneGenNoNL
		}
		r := sub.value(t).(rune)
		_, _ = w.WriteRune(maybeFoldCase(t.s, r, re.Flags))
	case syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText,
		syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		t.s.drawBits(0) // do nothing and hope that find() is enough
	case syntax.OpCapture:
		g.build(w, re.Sub[0], t)
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
		for repeat.more(t.s, regexpName(re.Sub[0])) {
			g.build(w, re.Sub[0], t)
		}
	case syntax.OpConcat:
		for _, sub := range re.Sub {
			g.build(w, sub, t)
		}
	case syntax.OpAlternate:
		ix := genIndex(t.s, len(re.Sub), true)
		g.build(w, re.Sub[ix], t)
	default:
		assertf(false, "invalid regexp op %v", re.Op)
	}

	t.s.endGroup(i, false)
}

func maybeFoldCase(s bitStream, r rune, flags syntax.Flags) rune {
	n := uint64(0)
	if flags&syntax.FoldCase != 0 {
		n, _, _ = genUintN(s, 4, false)
	}

	for i := 0; i < int(n); i++ {
		r = unicode.SimpleFold(r)
	}

	return r
}

func expandRangeTable(t *unicode.RangeTable, key interface{}) []rune {
	cached, ok := expandedTables.Load(key)
	if ok {
		return cached.([]rune)
	}

	var ret []rune
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
	expandedTables.Store(key, ret)

	return ret
}

func compileRegexp(expr string) (compiledRegexp, error) {
	cached, ok := compiledRegexps.Load(expr)
	if ok {
		return cached.(compiledRegexp), nil
	}

	syn, err := syntax.Parse(expr, syntax.Perl)
	if err != nil {
		return compiledRegexp{}, fmt.Errorf("failed to parse regexp %q: %v", expr, err)
	}

	re, err := regexp.Compile(expr)
	if err != nil {
		return compiledRegexp{}, fmt.Errorf("failed to compile regexp %q: %v", expr, err)
	}

	ret := compiledRegexp{syn, re}
	compiledRegexps.Store(expr, ret)

	return ret, nil
}

func regexpName(re *syntax.Regexp) string {
	cached, ok := regexpNames.Load(re)
	if ok {
		return cached.(string)
	}

	s := re.String()
	regexpNames.Store(re, s)

	return s
}

func charClassGen(re *syntax.Regexp) *Generator {
	cached, ok := charClassGens.Load(regexpName(re))
	if ok {
		return cached.(*Generator)
	}

	t := &unicode.RangeTable{}
	for i := 0; i < len(re.Rune); i += 2 {
		t.R32 = append(t.R32, unicode.Range32{
			Lo:     uint32(re.Rune[i]),
			Hi:     uint32(re.Rune[i+1]),
			Stride: 1,
		})
	}

	g := newGenerator(&runeGen{
		die:    newLoadedDie([]int{1}),
		tables: [][]rune{expandRangeTable(t, regexpName(re))},
	})
	charClassGens.Store(regexpName(re), g)

	return g
}
