package gotimer

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	m.Run()
}

func Test_timer(t *testing.T) {
	timer := NewTimer()
	job := timer.AddFunc(time.Second, func() {
		t.Log("job1", time.Now().Local())
	})
	time.Sleep(5 * time.Second)
	timer.Remove(job)
	job2 := timer.AddFunc(time.Second, func() {
		t.Log("job2", time.Now().Local())
	})
	t.Log("job2", job2)
	time.Sleep(5 * time.Second)
	timer.Clear()
	time.Sleep(time.Minute)
}
