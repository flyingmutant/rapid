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
type typeEnum byte

const (
	UNKNOWN typeEnum = iota
	UINT
	INT
	FLOAT
)

type numCounter struct {
	ints   []int64
	uints  []uint64
	floats []float64
	t      typeEnum
}

// stats maps labels to counters
type stats map[string]counter
type nStats map[string]*numCounter
type counterPair struct {
	frequency int
	event     string
}

// Event records an event for test `t` and
// stores the event for calculating statistics.
//
// Recording events and printing a their statistic is a tool for
// analysing test data generation. It helps to understand if
// your (custom) generators produce values in the expected range.
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

// NumericEvent records a numeric event for given label
// for statistic calculation of minimum, median, arithmetic mean
// and maximum values for that label. Allowed numeric values are
// all signed and unsigned integers and floats. Howver, for each
// label, all numeric values must be consistent, ie. have the the
// same type. Otherwise, a run-time panic will occur.
//
// Recording events and printing a their statistic is a tool for
// analysing test data generation. It helps to understand if
// your (custom) generators produce values in the expected range.
//
// Each event has a label and an event value. To see the statistics,
// run the tests with `go test -v`.
func NumericEvent(t *T, label string, numValue value) {
	valType := UNKNOWN
	var uintVal uint64
	var intVal int64
	var floatVal float64
	if t.tb != nil {
		t.tb.Helper()
	}
	switch x := numValue.(type) {
	case uint8:
		valType = UINT
		uintVal = uint64(x)
	case uint16:
		valType = UINT
		uintVal = uint64(x)
	case uint32:
		valType = UINT
		uintVal = uint64(x)
	case uint64:
		valType = UINT
		uintVal = uint64(x)
	case uint:
		valType = UINT
		uintVal = uint64(x)
	case int8:
		valType = INT
		intVal = int64(x)
	case int16:
		valType = INT
		intVal = int64(x)
	case int32:
		valType = INT
		intVal = int64(x)
	case int64:
		valType = INT
		intVal = int64(x)
	case int:
		valType = INT
		intVal = int64(x)
	case float32:
		valType = FLOAT
		floatVal = float64(x)
	case float64:
		valType = FLOAT
		floatVal = float64(x)
	default:
		t.Fatalf("numeric event is not a numeric value")
		return
	}
	t.statMux.Lock()
	defer t.statMux.Unlock()
	if t.numStats == nil {
		t.numStats = make(nStats)
	}
	c, found := t.numStats[label]
	if !found {
		newCounter := &numCounter{
			ints:   []int64{},
			uints:  []uint64{},
			floats: []float64{},
			t:      valType,
		}
		t.numStats[label] = newCounter
		c = newCounter
	}
	if valType != c.t {
		t.Fatalf("Type of numeric event does not match. Expected: %#v, was %#v", c.t, valType)
	}
	switch valType {
	case UINT:
		c.uints = append(c.uints, uintVal)
	case INT:
		c.ints = append(c.ints, intVal)
	case FLOAT:
		c.floats = append(c.floats, floatVal)
	}
}

func minMaxMeanInt(values []int64) (min, max, median int64, mean float64) {
	min, max, mean = values[0], values[0], float64(values[0])
	sum := mean
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
		if values[i] > max {
			max = values[i]
		}
		sum += float64(values[i])
	}
	mean = sum / float64(len(values))
	sort.Sort(int64Slice(values))
	if len(values)%2 == 1 {
		median = values[(len(values)-1)/2]
	} else {
		n := len(values) / 2
		median = (values[n] + values[n-1]) / 2
	}
	return
}

func minMaxMeanUint(values []uint64) (min, max, median uint64, mean float64) {
	min, max, mean = values[0], values[0], float64(values[0])
	sum := mean
	// log.Printf("rapid mean: values[%d] = %v", 0, values[0])
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
		if values[i] > max {
			max = values[i]
		}
		// log.Printf("rapid mean: values[%d] = %v", i, values[i])
		sum += float64(values[i])
	}
	mean = sum / float64(len(values))
	sort.Sort(uint64Slice(values))
	if len(values)%2 == 1 {
		median = values[(len(values)-1)/2]
	} else {
		n := len(values) / 2
		median = (values[n] + values[n-1]) / 2
	}
	return
}

func minMaxMeanFloat(values []float64) (min, max, median, mean float64) {
	min, max, mean = values[0], values[0], values[0]
	sum := mean
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
		if values[i] > max {
			max = values[i]
		}
		sum += values[i]
	}
	mean = sum / float64(len(values))
	sort.Float64s(values)
	if len(values)%2 == 1 {
		median = values[(len(values)-1)/2]
	} else {
		n := len(values) / 2
		median = (values[n] + values[n-1]) / 2
	}
	return
}

// printStats logs a table of events and their relative frequency.
func printStats(t *T) {
	t.statMux.Lock()
	defer t.statMux.Unlock()
	if t.tbLog && t.tb != nil {
		t.tb.Helper()
	}
	if len(t.allstats) > 0 {
		log.Printf("[rapid] Statistics for %s\n", t.Name())
		for label := range t.allstats {
			log.Printf("[rapid] Events with label %s", label)
			s := t.allstats[label]
			events := make([]counterPair, 0)
			sum := 0
			count := 0
			for ev := range s {
				sum += s[ev]
				count++
				events = append(events, counterPair{event: ev, frequency: s[ev]})
			}
			log.Printf("[rapid] Total of %d different events\n", count)
			// we sort twice to sort same frequency alphabetically
			sort.Slice(events, func(i, j int) bool { return events[i].event < events[j].event })
			sort.SliceStable(events, func(i, j int) bool { return events[i].frequency > events[j].frequency })
			for _, ev := range events {
				log.Printf("[rapid]   %s: %d (%f %%)\n", ev.event, ev.frequency, float32(ev.frequency)/float32(sum)*100.0)
			}
		}
	}
	if len(t.numStats) > 0 {
		log.Printf("[rapid] Numerical Statistics for %s\n", t.Name())
		for label := range t.numStats {
			log.Printf("[rapid] Numerical events with label %s", label)
			c := t.numStats[label]
			switch c.t {
			case UINT:
				min, max, median, mean := minMaxMeanUint(c.uints)
				log.Printf("[rapid]   Min:    %d\n", min)
				log.Printf("[rapid]   Median: %d\n", median)
				log.Printf("[rapid]   Mean:   %f\n", mean)
				log.Printf("[rapid]   Max:    %d\n", max)
			case INT:
				min, max, median, mean := minMaxMeanInt(c.ints)
				log.Printf("[rapid]   Min:    %d\n", min)
				log.Printf("[rapid]   Median: %d\n", median)
				log.Printf("[rapid]   Mean:   %f\n", mean)
				log.Printf("[rapid]   Max:    %d\n", max)
			case FLOAT:
				min, max, median, mean := minMaxMeanFloat(c.floats)
				log.Printf("[rapid]   Min:    %f\n", min)
				log.Printf("[rapid]   Median: %f\n", median)
				log.Printf("[rapid]   Mean:   %f\n", mean)
				log.Printf("[rapid]   Max:    %f\n", max)
			}
		}
	}
}

// implementing the sorting interfaces for int64 and uint64 slices
type uint64Slice []uint64
type int64Slice []int64

func (x uint64Slice) Len() int           { return len(x) }
func (x uint64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x uint64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x int64Slice) Len() int            { return len(x) }
func (x int64Slice) Less(i, j int) bool  { return x[i] < x[j] }
func (x int64Slice) Swap(i, j int)       { x[i], x[j] = x[j], x[i] }
