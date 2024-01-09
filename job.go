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
				slog.Info("stop job", slog.String("job", j.id))
				return
			case <-j.ticker.C:
				slog.Debug("do job", slog.String("job", j.id))
				j.do(cmd)
			}
		}
	}()
}

func (j *Job) do(cmd func()) {
	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case error:
				slog.Error(fmt.Sprintf("job panic:%v", v), slog.String("job", j.id))
			default:
				slog.Error(fmt.Sprintf("job panic:%v", v), slog.String("job", j.id))
			}
		}
	}()
	cmd()
}

func (j *Job) stop() {
	if j.cancel != nil {
		j.cancel()
	}
	if j.ticker != nil {
		j.ticker.Stop()
	}
}
