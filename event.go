//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
package rapid

import (
	"sort"
)

type counter map[string]int
type stats map[string]counter

type counterPair struct {
	frequency int
	event     string
}

// global variable for statistics
var all_stats stats = make(stats)

// Event records an event for test `t` and
// stores the event for calculating statistics.
//
// Recording events and printing a their statistic is a tool for
// analysing test data generations. It helps to understand if
// your customer generators produce value in the expected range.
func Event(t TB, event string) {
	t.Helper()
	c, found := all_stats[t.Name()]
	if !found {
		c = make(counter)
		all_stats[t.Name()] = c
	}
	c[event]++
}

// PrintStats logs a table of events and their relative frequency.
// To see these statistics, run the tests with `go test -v`
func PrintStats(t TB) {
	t.Helper()
	s, found := all_stats[t.Name()]
	if !found {
		t.Logf("No events stored for test %s", t)
		return
	}
	events := make([]counterPair, 0)
	sum := 0
	count := 0
	for ev := range s {
		sum += s[ev]
		count++
		events = append(events, counterPair{event: ev, frequency: s[ev]})
	}
	t.Logf("Statistics for %s\n", t.Name())
	t.Logf("Total of %d different events\n", count)
	t.Logf(" =====> Sorted by frequency:")
	sort.Slice(events, func(i, j int) bool { return events[i].frequency > events[j].frequency })
	for _, ev := range events {
		t.Logf("%s: %d (%f %%)\n", ev.event, ev.frequency, float32(ev.frequency)/float32(sum)*100.0)
	}
}
