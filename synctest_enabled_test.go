//go:build go1.25

package rapid_test

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"pgregory.net/rapid"
)

func TestSyncTest(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		var cleaned atomic.Bool

		rapid.SyncTest(rt, func(inner *rapid.T) {
			inner.Cleanup(func() {
				cleaned.Store(true)
			})

			const sleep = 2 * time.Second
			start := time.Now()
			time.Sleep(sleep)
			if got := time.Since(start); got != sleep {
				inner.Fatalf("virtual time advanced by %v, want %v", got, sleep)
			}

			done := make(chan struct{})
			go func() { close(done) }()
			synctest.Wait()
			select {
			case <-done:
			default:
				inner.Fatalf("goroutine did not finish inside synctest bubble")
			}
		})

		if !cleaned.Load() {
			rt.Fatalf("cleanup registered inside SyncTest did not run")
		}
	})
}
