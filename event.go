//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
package rapid

import (
	"log"
	"sort"
)

// counter maps a (stringified) event to a frequency counter
type counter map[string]int

// stats maps labels to counters
type stats map[string]counter
type counterPair struct {
	frequency int
	event     string
}

// Event records an event for test `t` and
// stores the event for calculating statistics.
//
// Recording events and printing a their statistic is a tool for
// analysing test data generations. It helps to understand if
// your customer generators produce value in the expected range.
//
// Each event has a label and an event value. To see the statistics,
// run the tests with `go test -v`.
//
func Event(t *T, label string, value string) {
	if t.tb != nil {
		t.tb.Helper()
	}
	t.statMux.Lock()
	defer t.statMux.Unlock()
	if t.allstats == nil {
		t.allstats = make(stats)
	}
	c, found := t.allstats[label]
	if !found {
		c = make(counter)
		t.allstats[label] = c
	}
	c[value]++
}

// printStats logs a table of events and their relative frequency.
func printStats(t *T) {
	t.statMux.Lock()
	defer t.statMux.Unlock()
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	if len(t.allstats) > 0 {
		log.Printf("Statistics for %s\n", t.Name())
		for label := range t.allstats {
			log.Printf("Events with label %s", label)
			s := t.allstats[label]
			events := make([]counterPair, 0)
			sum := 0
			count := 0
			for ev := range s {
				sum += s[ev]
				count++
				events = append(events, counterPair{event: ev, frequency: s[ev]})
			}
			log.Printf("Total of %d different events\n", count)
			// we sort twice to sort same frequency alphabetically
			sort.Slice(events, func(i, j int) bool { return events[i].event < events[j].event })
			sort.SliceStable(events, func(i, j int) bool { return events[i].frequency > events[j].frequency })
			for _, ev := range events {
				log.Printf("  %s: %d (%f %%)\n", ev.event, ev.frequency, float32(ev.frequency)/float32(sum)*100.0)
			}
		}
	}
}
