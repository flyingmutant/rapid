// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import (
	"encoding/binary"
	"fmt"
	"math/bits"
	"os"
	"strings"
	"time"
)

const shrinkTimeLimit = 30 * time.Second

func shrink(tb limitedTB, rec recordedBits, err *testError, prop func(*T)) ([]uint64, *testError) {
	rec.prune()

	s := &shrinker{
		tb:    tb,
		rec:   rec,
		err:   err,
		prop:  prop,
		cache: map[string]struct{}{},
	}

	buf, err := s.shrink()

	if *debugvis {
		name := fmt.Sprintf("vis-%v.html", tb.Name())
		f, err := os.Create(name)
		if err != nil {
			tb.Logf("failed to create debugvis file %v: %v", name, err)
		} else {
			defer f.Close()

			if err = visWriteHTML(f, tb.Name(), s.visBits); err != nil {
				tb.Logf("failed to write debugvis file %v: %v", name, err)
			}
		}
	}

	return buf, err
}

type shrinker struct {
	tb      limitedTB
	rec     recordedBits
	err     *testError
	prop    func(*T)
	visBits []recordedBits
	tries   int
	shrinks int
	cache   map[string]struct{}
	hits    int
}

func (s *shrinker) debugf(format string, args ...interface{}) {
	if *debug {
		s.tb.Helper()
		s.tb.Logf("[shrink] "+format, args...)
	}
}

func (s *shrinker) shrink() (buf []uint64, err *testError) {
	defer func() {
		if r := recover(); r != nil {
			buf, err = s.rec.data, r.(*testError)
		}
	}()

	i := 0
	shrinks := -1
	start := time.Now()
	for ; s.shrinks > shrinks && time.Since(start) < shrinkTimeLimit; i++ {
		shrinks = s.shrinks

		s.debugf("round %v start", i)
		s.removeGroups()
		s.minimizeBlocks()
	}
	s.debugf("done, %v rounds total (%v tries, %v shrinks, %v cache hits)", i, s.tries, s.shrinks, s.hits)

	return s.rec.data, s.err
}

func (s *shrinker) removeGroups() {
	for i := 0; i < len(s.rec.groups); i++ {
		g := s.rec.groups[i]
		if !g.standalone || g.end < 0 {
			continue
		}

		if s.accept(without(s.rec.data, g), "remove group %q at %v: [%v, %v)", g.label, i, g.begin, g.end) {
			i--
		}
	}
}

func (s *shrinker) minimizeBlocks() {
	for i := 0; i < len(s.rec.data); i++ {
		minimize(s.rec.data[i], func(u uint64) bool {
			buf := append([]uint64(nil), s.rec.data...)
			buf[i] = u
			return s.accept(buf, "minimize block %v: %v to %v", i, s.rec.data[i], u)
		})
	}
}

func (s *shrinker) accept(buf []uint64, format string, args ...interface{}) bool {
	if compareData(buf, s.rec.data) >= 0 {
		return false
	}
	bufStr := dataStr(buf)
	if _, ok := s.cache[bufStr]; ok {
		s.hits++
		return false
	}

	s.debugf("trying to reproduce the failure with a smaller test case: "+format, args...)
	s.tries++
	s1 := newBufBitStream(buf, false)
	err1 := checkOnce(newT(s.tb, s1, *debug), s.prop)
	if traceback(err1) != traceback(s.err) {
		s.cache[bufStr] = struct{}{}
		return false
	}

	s.debugf("trying to reproduce the failure")
	s.err = err1
	s2 := newBufBitStream(buf, true)
	err2 := checkOnce(newT(s.tb, s2, *debug), s.prop)
	s.rec = s2.recordedBits
	s.rec.prune()
	assert(compareData(s.rec.data, buf) <= 0)
	if *debugvis {
		s.visBits = append(s.visBits, s.rec)
	}
	if !sameError(err1, err2) {
		panic(err2)
	}
	s.shrinks++

	return true
}

func minimize(u uint64, cond func(uint64) bool) uint64 {
	if u == 0 {
		return 0
	}
	for i := uint64(0); i < u && i < small; i++ {
		if cond(i) {
			return i
		}
	}
	if u <= small {
		return u
	}

	m := &minimizer{best: u, cond: cond}

	m.rShift()
	m.unsetBits()
	m.sortBits()
	m.binSearch()

	return m.best
}

type minimizer struct {
	best uint64
	cond func(uint64) bool
}

func (m *minimizer) accept(u uint64) bool {
	if u >= m.best || !m.cond(u) {
		return false
	}
	m.best = u
	return true
}

func (m *minimizer) rShift() {
	for m.accept(m.best >> 1) {
	}
}

func (m *minimizer) unsetBits() {
	size := bits.Len64(m.best)

	for i := size - 1; i >= 0; i-- {
		m.accept(m.best ^ 1<<uint(i))
	}
}

func (m *minimizer) sortBits() {
	size := bits.Len64(m.best)

	for i := size - 1; i >= 0; i-- {
		h := uint64(1 << uint(i))
		if m.best&h != 0 {
			for j := 0; j < i; j++ {
				l := uint64(1 << uint(j))
				if m.best&l == 0 {
					if m.accept(m.best ^ (l | h)) {
						break
					}
				}
			}
		}
	}
}

func (m *minimizer) binSearch() {
	if !m.accept(m.best - 1) {
		return
	}

	i := uint64(0)
	j := m.best
	for i < j {
		h := i + (j-i)/2
		if m.accept(h) {
			j = h
		} else {
			i = h + 1
		}
	}
}

func without(data []uint64, groups ...groupInfo) []uint64 {
	buf := append([]uint64(nil), data...)

	for i := len(groups) - 1; i >= 0; i-- {
		g := groups[i]
		buf = append(buf[:g.begin], buf[g.end:]...)
	}

	return buf
}

func dataStr(data []uint64) string {
	b := &strings.Builder{}
	err := binary.Write(b, binary.BigEndian, data)
	assert(err == nil)
	return b.String()
}

func compareData(a []uint64, b []uint64) int {
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}

	for i := range a {
		if a[i] < b[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}

	return 0
}
