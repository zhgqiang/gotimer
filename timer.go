package gotimer

import (
	"sync"
	"time"
)

type Timer struct {
	lock sync.Mutex
	jobs map[string]*Job
}

func NewTimer() *Timer {
	t := &Timer{
		jobs: map[string]*Job{},
	}
	return t
}

func (t *Timer) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()
	for _, job := range t.jobs {
		job.stop()
	}
	clear(t.jobs)
}

func (t *Timer) AddFunc(duration time.Duration, cmd func()) string {
	job := NewJob(duration)
	t.saveTicker(job.id, job)
	job.start(cmd)
	return job.id
}

func (t *Timer) saveTicker(id string, job *Job) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.jobs[id] = job
}

func (t *Timer) Remove(id string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	job, ok := t.jobs[id]
	if ok {
		job.stop()
		delete(t.jobs, id)
	}
}
