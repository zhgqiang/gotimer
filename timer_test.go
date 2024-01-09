package gotimer

import (
	"fmt"
	"github.com/pkg/errors"
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

func Test_timer2(t *testing.T) {
	timer := NewTimer()
	job := timer.AddFunc(time.Second, func() {
		t.Log("job1", time.Now().Local())
	})
	t.Log("job1", job)
	job2 := timer.AddFunc(time.Second, func() {
		t.Log("job2", time.Now().Local())
	})
	t.Log("job2", job2)
	job3 := timer.AddFunc(time.Second, func() {
		err := fmt.Errorf("job3 err1")
		err = errors.Wrap(err, "err2")
		panic(err)
	})
	t.Log("job3", job3)
	time.Sleep(time.Second * 10)
	timer.Clear()
}
