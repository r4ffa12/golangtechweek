package wokerpool

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Job interface{}

type Result interface{}

type ProcessFunc func(ctx context.Context, job Job) Result

type WorkerPool interface {
	Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error)
	Stop() error
	IsRunning() bool
}

type State int

const (
	StartIdle State = iota
	StateRunning
	StateStopped
)

type Config struct {
	WorkerCount int
	Logger      *slog.Logger
}

func DefaultConfig() Config {
	return Config{
		WorkerCount: 1,
		Logger:      slog.Default(),
	}
}

type workerPool struct {
	WorkerCount int
	ProcessFunc ProcessFunc
	logger      *slog.Logger
	state       State
	stateMutex  sync.Mutex
	stopCh      chan struct{}
	stopWg      sync.WaitGroup
}

func New(ProcessFunc ProcessFunc, config Config) *workerPool {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 1
	}
	if config.Logger == nil {
		config.Logger = slog.Default()
	}
	return &workerPool{
		ProcessFunc: ProcessFunc,
		WorkerCount: config.WorkerCount,
		stopCh:      make(chan struct{}),
		state:       StartIdle,
		logger:      config.Logger,
	}
}

func (wp *workerPool) Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error) {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()

	if wp.state != StartIdle {
		return nil, fmt.Errorf("worker pool is not in idle state")
	}

	resultCh := make(chan Result)
	wp.state = StateRunning
	wp.stopCh = make(chan struct{})

	wp.stopWg.Add(wp.WorkerCount)

	for i := 0; i < wp.WorkerCount; i++ {
		go wp.worker(ctx, i, inputCh, resultCh)
	}

	go func() {
		wp.stopWg.Wait()
		close(resultCh)

		wp.stateMutex.Lock()
		wp.state = StartIdle
		wp.stateMutex.Unlock()
	}()

	return resultCh, nil

}

func (wp *workerPool) Stop() error {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()

	if wp.state != StateRunning {
		return fmt.Errorf("worker pool is not in running state")
	}

	wp.state = StateStopped
	close(wp.stopCh)

	wp.stopWg.Wait()

	wp.state = StartIdle
	return nil
}

func (wp *workerPool) IsRunning() bool {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()
	return wp.state == StateRunning
}

func (wp *workerPool) worker(ctx context.Context, id int, inputCh <-chan Job, resultCh chan<- Result) {
	wp.logger.Info("worker started", "worker_id", id)

	for {
		select {
		case <-wp.stopCh:
			wp.logger.Info("worker interrompido", "worker_id", id)
			return
		case <-ctx.Done():
			wp.logger.Info("Contexto cancelado, interrompendo worker", "worker_id", id)
			return
		case job, ok := <-inputCh:
			if !ok {
				wp.logger.Info("inputCh closed, stopping worker", "worker_id", id)
				return
			}

			result := wp.ProcessFunc(ctx, job)
			select {
			case resultCh <- result:
			case <-wp.stopCh:
				wp.logger.Info("worker interrompido", "worker_id", id)
				return

			case <-ctx.Done():
				wp.logger.Info("Contexto cancelado, interrompendo worker", "worker_id", id)
			}
		}
	}
}
