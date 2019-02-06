// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rapid_test

import (
	"testing"

	"github.com/flyingmutant/rapid"
)

// Queue implements integer queue with a fixed maximum size.
type Queue struct {
	buf []int
	in  int
	out int
}

func NewQueue(n int) *Queue {
	return &Queue{
		buf: make([]int, n+1),
	}
}

// Precondition: Size() > 0.
func (q *Queue) Get() int {
	i := q.buf[q.out]
	q.out = (q.out + 1) % len(q.buf)
	return i
}

// Precondition: Size() < n.
func (q *Queue) Put(i int) {
	q.buf[q.in] = i
	q.in = (q.in + 1) % len(q.buf)
}

func (q *Queue) Size() int {
	return (q.in - q.out) % len(q.buf)
}

// queueMachine is a description of a rapid state machine for testing Queue
type queueMachine struct {
	q     *Queue // queue being tested
	n     int    // maximum queue size
	state []int  // model of the queue
}

// Init is an action for initializing  a queueMachine instance.
func (m *queueMachine) Init() func(*rapid.T) {
	return rapid.Bind(m.init, rapid.IntsRange(0, 1000))
}

// Get is a conditional action which removes an item from the queue.
func (m *queueMachine) Get() func(*rapid.T) {
	return rapid.BindIf(m.q.Size() > 0, m.get)
}

// Put is a conditional action which adds an items to the queue.
func (m *queueMachine) Put() func(*rapid.T) {
	return rapid.BindIf(m.q.Size() < m.n, m.put, rapid.Ints())
}

// Check verifies that all required invariants hold.
func (m *queueMachine) Check(t *rapid.T) {
	if m.q.Size() != len(m.state) {
		t.Fatalf("queue size mismatch: %v vs expected %v", m.q.Size(), len(m.state))
	}
}

func (m *queueMachine) init(t *rapid.T, n int) {
	m.q = NewQueue(n)
	m.n = n
}

func (m *queueMachine) get(t *rapid.T) {
	i := m.q.Get()
	if i != m.state[0] {
		t.Fatalf("got invalid value: %v vs expected %v", i, m.state[0])
	}
	m.state = m.state[1:]
}

func (m *queueMachine) put(t *rapid.T, i int) {
	m.q.Put(i)
	m.state = append(m.state, i)
}

// Rename to TestQueue to make an actual (failing) test.
func Example_queue(t *testing.T) {
	rapid.Check(t, rapid.StateMachine(&queueMachine{}))
}
