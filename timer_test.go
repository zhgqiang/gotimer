package gotimer

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func waitFor(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}

	t.Fatal("condition not satisfied before timeout")
}

func TestTimerAddFuncRunsPeriodically(t *testing.T) {
	timer := NewTimer()
	t.Cleanup(timer.Clear)

	var count atomic.Int32

	jobID := timer.AddFunc(10*time.Millisecond, func() {
		count.Add(1)
	})

	if jobID == "" {
		t.Fatal("expected non-empty job id")
	}

	waitFor(t, 200*time.Millisecond, func() bool {
		return count.Load() >= 3
	})
}

func TestTimerRemoveStopsJob(t *testing.T) {
	timer := NewTimer()
	var count atomic.Int32

	jobID := timer.AddFunc(10*time.Millisecond, func() {
		count.Add(1)
	})

	waitFor(t, 200*time.Millisecond, func() bool {
		return count.Load() >= 2
	})

	timer.Remove(jobID)
	current := count.Load()

	time.Sleep(50 * time.Millisecond)

	if after := count.Load(); after != current {
		t.Fatalf("expected job to stop after Remove, before=%d after=%d", current, after)
	}
}

func TestTimerClearStopsAllJobsAndClearsMap(t *testing.T) {
	timer := NewTimer()
	var count1 atomic.Int32
	var count2 atomic.Int32

	timer.AddFunc(10*time.Millisecond, func() {
		count1.Add(1)
	})
	timer.AddFunc(10*time.Millisecond, func() {
		count2.Add(1)
	})

	waitFor(t, 200*time.Millisecond, func() bool {
		return count1.Load() >= 1 && count2.Load() >= 1
	})

	timer.Clear()

	if got := len(timer.jobs); got != 0 {
		t.Fatalf("expected no jobs after Clear, got %d", got)
	}

	before1 := count1.Load()
	before2 := count2.Load()

	time.Sleep(50 * time.Millisecond)

	if after1 := count1.Load(); after1 != before1 {
		t.Fatalf("expected first job to stop after Clear, before=%d after=%d", before1, after1)
	}
	if after2 := count2.Load(); after2 != before2 {
		t.Fatalf("expected second job to stop after Clear, before=%d after=%d", before2, after2)
	}
}

func TestTimerRecoversFromPanic(t *testing.T) {
	timer := NewTimer()
	t.Cleanup(timer.Clear)

	var count atomic.Int32

	timer.AddFunc(10*time.Millisecond, func() {
		count.Add(1)
		err := fmt.Errorf("job error")
		panic(errors.Wrap(err, "wrapped"))
	})

	waitFor(t, 200*time.Millisecond, func() bool {
		return count.Load() >= 3
	})
}

func TestTimerRemoveUnknownIDDoesNothing(t *testing.T) {
	timer := NewTimer()
	t.Cleanup(timer.Clear)

	timer.AddFunc(10*time.Millisecond, func() {})

	before := len(timer.jobs)
	timer.Remove("not-exists")
	after := len(timer.jobs)

	if after != before {
		t.Fatalf("expected removing unknown id to keep jobs unchanged, before=%d after=%d", before, after)
	}
}

func TestTimerClearIsIdempotent(t *testing.T) {
	timer := NewTimer()

	timer.AddFunc(10*time.Millisecond, func() {})
	timer.Clear()
	timer.Clear()

	if got := len(timer.jobs); got != 0 {
		t.Fatalf("expected no jobs after repeated Clear, got %d", got)
	}
}

func TestTimerRemoveIsIdempotent(t *testing.T) {
	timer := NewTimer()

	jobID := timer.AddFunc(10*time.Millisecond, func() {})
	timer.Remove(jobID)
	timer.Remove(jobID)

	if got := len(timer.jobs); got != 0 {
		t.Fatalf("expected no jobs after repeated Remove, got %d", got)
	}
}

func TestTimerConcurrentAddAndRemove(t *testing.T) {
	timer := NewTimer()
	t.Cleanup(timer.Clear)

	const workers = 10
	const jobsPerWorker = 20

	ids := make(chan string, workers*jobsPerWorker)
	var addWG sync.WaitGroup

	for i := 0; i < workers; i++ {
		addWG.Add(1)
		go func() {
			defer addWG.Done()
			for j := 0; j < jobsPerWorker; j++ {
				id := timer.AddFunc(time.Hour, func() {})
				ids <- id
			}
		}()
	}

	addWG.Wait()
	close(ids)

	if got := len(timer.jobs); got != workers*jobsPerWorker {
		t.Fatalf("expected %d jobs after concurrent add, got %d", workers*jobsPerWorker, got)
	}

	var removeWG sync.WaitGroup
	for id := range ids {
		removeWG.Add(1)
		go func(jobID string) {
			defer removeWG.Done()
			timer.Remove(jobID)
		}(id)
	}

	removeWG.Wait()

	if got := len(timer.jobs); got != 0 {
		t.Fatalf("expected no jobs after concurrent remove, got %d", got)
	}
}

func BenchmarkTimerAddFunc(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		timer := NewTimer()
		_ = timer.AddFunc(time.Hour, func() {})
		timer.Clear()
	}
}

func BenchmarkTimerRemove(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		timer := NewTimer()
		jobID := timer.AddFunc(time.Hour, func() {})
		timer.Remove(jobID)
	}
}

func BenchmarkJobDo(b *testing.B) {
	job := NewJob(time.Hour)
	defer job.stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		job.do(func() {})
	}
}
