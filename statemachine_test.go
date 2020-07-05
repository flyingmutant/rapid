// Copyright 2019 Gregory Petrosyan <gregory.petrosyan@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid

import "testing"

// https://github.com/leanovate/gopter/blob/master/commands/example_circularqueue_test.go
var gopterBug = false

// https://godoc.org/github.com/leanovate/gopter/commands#example-package--BuggyCounter
type buggyCounter struct {
	n int
}

func (c *buggyCounter) Get() int {
	return c.n
}

func (c *buggyCounter) Inc() {
	c.n++
}

func (c *buggyCounter) Dec() {
	if c.n > 3 {
		c.n -= 2
	} else {
		c.n--
	}
}

func (c *buggyCounter) Reset() {
	c.n = 0
}

type counterMachine struct {
	c    buggyCounter
	incs int
	decs int
}

func (m *counterMachine) Inc(_ *T) {
	m.c.Inc()
	m.incs++
}

func (m *counterMachine) Dec(_ *T) {
	m.c.Dec()
	m.decs++
}

func (m *counterMachine) Reset(_ *T) {
	m.c.Reset()
	m.incs = 0
	m.decs = 0
}

func (m *counterMachine) Check(t *T) {
	if m.c.Get() != m.incs-m.decs {
		t.Fatalf("counter value is %v with %v incs and %v decs", m.c.Get(), m.incs, m.decs)
	}
}

func TestStateMachine_Counter(t *testing.T) {
	t.Parallel()

	checkShrink(t, Run(&counterMachine{}),
		"Inc", "Inc", "Inc", "Inc",
		"Dec",
	)
}

type haltingMachine struct {
	a []int
	b []int
	c []int
}

func (m *haltingMachine) Check(t *T) {
	if len(m.a) > 3 || len(m.b) > 3 || len(m.c) > 3 {
		t.Fatalf("too many elements: %v, %v, %v", len(m.a), len(m.b), len(m.c))
	}
}

func (m *haltingMachine) A(t *T) {
	if len(m.a) == 3 {
		t.SkipNow()
	}

	m.a = append(m.a, Int().Draw(t, "a").(int))
}

func (m *haltingMachine) B(t *T) {
	if len(m.b) == 3 {
		t.SkipNow()
	}

	m.b = append(m.b, Int().Draw(t, "b").(int))
}

func (m *haltingMachine) C(t *T) {
	if len(m.c) == 3 {
		t.SkipNow()
	}

	m.c = append(m.c, Int().Draw(t, "c").(int))
}

func TestStateMachine_Halting(t *testing.T) {
	t.Parallel()

	a := []value{"A", 0, "A", 0, "A", 0}
	for i := 0; i < 100; i++ {
		a = append(a, "A") // TODO proper shrinking of "stuck" state machines
	}

	checkShrink(t, Run(&haltingMachine{}), a...)
}

// https://www.cs.tufts.edu/~nr/cs257/archive/john-hughes/quviq-testing.pdf
type buggyQueue struct {
	buf []int
	in  int
	out int
}

func newBuggyQueue(size int) *buggyQueue {
	return &buggyQueue{
		buf: make([]int, size+1),
	}
}

func (q *buggyQueue) Get() int {
	n := q.buf[q.out]
	q.out = (q.out + 1) % len(q.buf)
	return n
}

func (q *buggyQueue) Put(i int) {
	if gopterBug && q.in == 4 && i > 0 {
		q.buf[len(q.buf)-1] *= i
	}

	q.buf[q.in] = i
	q.in = (q.in + 1) % len(q.buf)
}

func (q *buggyQueue) Size() int {
	if gopterBug {
		return (q.in - q.out + len(q.buf)) % len(q.buf)
	} else {
		return (q.in - q.out) % len(q.buf)
	}
}

type queueMachine struct {
	q     *buggyQueue
	state []int
	size  int
}

func (m *queueMachine) Init(t *T) {
	size := IntRange(1, 1000).Draw(t, "size").(int)
	m.q = newBuggyQueue(size)
	m.size = size
}

func (m *queueMachine) Get(t *T) {
	if m.q.Size() == 0 {
		t.Skip("queue empty")
	}

	n := m.q.Get()
	if n != m.state[0] {
		t.Fatalf("got invalid value: %v vs expected %v", n, m.state[0])
	}
	m.state = m.state[1:]
}

func (m *queueMachine) Put(t *T) {
	if m.q.Size() == m.size {
		t.Skip("queue full")
	}

	n := Int().Draw(t, "n").(int)
	m.q.Put(n)
	m.state = append(m.state, n)
}

func (m *queueMachine) Check(t *T) {
	if m.q.Size() != len(m.state) {
		t.Fatalf("queue size mismatch: %v vs expected %v", m.q.Size(), len(m.state))
	}
}

func TestStateMachine_Queue(t *testing.T) {
	t.Parallel()

	checkShrink(t, Run(&queueMachine{}),
		1,
		"Put", 0,
		"Get",
		"Put", 0,
	)
}

type garbageMachine struct {
	a []int
	b []int
}

func (m *garbageMachine) AddA(t *T) {
	if len(m.b) < 3 {
		t.Skip("too early")
	}

	n := Int().Draw(t, "a").(int)
	m.a = append(m.a, n)
}

func (m *garbageMachine) AddB(t *T) {
	n := Int().Draw(t, "b").(int)
	m.b = append(m.b, n)
}

func (m *garbageMachine) Whatever1(t *T) {
	b := Bool().Draw(t, "whatever 1/1").(bool)
	if b {
		t.Skip("arbitrary decision")
	}

	Float64().Draw(t, "whatever 1/2")
}

func (m *garbageMachine) Whatever2(t *T) {
	SliceOfDistinct(Int(), nil).Draw(t, "whatever 2")
}

func (m *garbageMachine) Whatever3(t *T) {
	OneOf(SliceOf(Byte()), MapOf(Int(), String())).Draw(t, "whatever 3")
}

func (m *garbageMachine) Check(t *T) {
	if len(m.a) > len(m.b) {
		t.Fatalf("`a` has outgrown `b`: %v vs %v", len(m.a), len(m.b))
	}
}

func TestStateMachine_DiscardGarbage(t *testing.T) {
	t.Parallel()

	checkShrink(t, Run(&garbageMachine{}),
		"AddB", 0,
		"AddB", 0,
		"AddB", 0,
		"AddA", 0,
		"AddA", 0,
		"AddA", 0,
		"AddA", 0,
	)
}

func BenchmarkCheckQueue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _ = doCheck(b, Run(&queueMachine{}))
	}
}
