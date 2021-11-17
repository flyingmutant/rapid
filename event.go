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

type event struct {
	label string
	value string
}
type done struct {
	result chan stats
}

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
	t.Helper()
	if t.evChan == nil {
		log.Printf("Creating the channels for test %s", t.Name())
		t.evChan = make(chan event)
		t.evDone = make(chan done)
		go eventRecorder(t.evChan, t.evDone)
	}
	ev := event{value: value, label: label}
	// log.Printf("Send the event %+v", ev)
	t.evChan <- ev
}

// eventRecorder is a goroutine that stores event for a test execution.
func eventRecorder(incomingEvent <-chan event, done <-chan done) {
	all_stats := make(stats)
	for {
		select {
		case ev := <-incomingEvent:
			c, found := all_stats[ev.label]
			if !found {
				c = make(counter)
				all_stats[ev.label] = c
			}
			c[ev.value]++
		case d := <-done:
			log.Printf("event Recorder: Done. Send the stats\n")
			d.result <- all_stats
			log.Printf("event Recorder: Done. Will return now\n")
			return
		}
	}
	// log.Printf("event Recorder: This shall never happen\n")

}

// printStats logs a table of events and their relative frequency.
func printStats(t *T) {
	// log.Printf("What about printing the stats for t = %+v", t)
	if t.evChan == nil || t.evDone == nil {
		return
	}
	log.Printf("Now we can print the stats")
	d := done{result: make(chan stats)}
	t.evDone <- d
	stats := <-d.result
	log.Printf("stats received")
	log.Printf("Statistics for %s\n", t.Name())
	for label := range stats {
		log.Printf("Events with label %s", label)
		s := stats[label]
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
			log.Printf("%s: %d (%f %%)\n", ev.event, ev.frequency, float32(ev.frequency)/float32(sum)*100.0)
		}
	}
	close(t.evChan)
	close(t.evDone)
}
