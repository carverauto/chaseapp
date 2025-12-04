package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Manager coordinates background workers and their lifecycle.
type Manager struct {
	logger  *slog.Logger
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	started bool
}

// NewManager creates a worker manager.
func NewManager(logger *slog.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

// Go starts a worker function in its own goroutine with context cancellation.
func (m *Manager) Go(name string, fn func(ctx context.Context)) {
	m.started = true
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error("worker panic", slog.String("worker", name), slog.Any("err", r))
			}
		}()
		fn(m.ctx)
	}()
}

// Stop cancels all workers and waits for completion.
func (m *Manager) Stop(ctx context.Context) {
	if m.cancel != nil {
		m.cancel()
	}
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		m.logger.Warn("timeout waiting for workers to stop")
	}
}

// RunInterval runs a task on a fixed interval until context cancellation.
func RunInterval(ctx context.Context, interval time.Duration, fn func(ctx context.Context)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately
	fn(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn(ctx)
		}
	}
}
