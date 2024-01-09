package gotimer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	id     string
	ctx    context.Context
	cancel context.CancelFunc
	ticker *time.Ticker
}

func NewJob(duration time.Duration) *Job {
	ctx, cancel := context.WithCancel(context.Background())
	id := uuid.New().String()
	ticker := time.NewTicker(duration)
	j := &Job{
		id:     id,
		ctx:    ctx,
		cancel: cancel,
		ticker: ticker,
	}
	return j
}

func (j *Job) start(cmd func()) {
	go func() {
		for {
			select {
			case <-j.ctx.Done():
				slog.Info(fmt.Sprintf("停止定时,id=%s", j.id))
				return
			case <-j.ticker.C:
				slog.Info(fmt.Sprintf("执行定时,id=%s", j.id))
				cmd()
			}
		}
	}()
}

func (j *Job) stop() {
	if j.cancel != nil {
		j.cancel()
	}
	if j.ticker != nil {
		j.ticker.Stop()
	}
}
