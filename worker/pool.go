package worker

import (
	"context"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"

	"github.com/site116/eventloader/config"
)

// Pool represents a worker pool that sends events to Kafka.
type Pool struct {
	cfg    config.Pool
	sender func(context.Context, int) error
}

// NewPool creates a new worker pool with the given configuration and sender function.
func NewPoll(cfg config.Pool, sender func(context.Context, int) error) *Pool {
	return &Pool{cfg: cfg, sender: sender}
}

// Run starts the worker pool and sends events to Kafka.
func (p *Pool) Run(ctx context.Context) {
	if p.sender == nil {
		log.Fatal("sender not set")
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, p.cfg.NumWorkers)
	var errorCounter atomic.Int32
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for i := 0; i < int(p.cfg.Batches); i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-sem
				wg.Done()
			}()
			if err := p.sender(ctx, i); err != nil {
				errorCounter.Add(1)
			}
			if errorCounter.Load() >= p.cfg.ErrorThreshold {
				cancel()
			}
		}()
	}
}
