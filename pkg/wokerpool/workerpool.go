package wokerpool

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Job representa um trabalho genérico a ser processado.
type Job interface{}

// Result representa o resultado do processamento de um trabalho.
type Result interface{}

// ProcessFunc define a função que processará os trabalhos recebidos.
type ProcessFunc func(ctx context.Context, job Job) Result

// WorkerPool define a interface para um pool de workers.
type WorkerPool interface {
	Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error)
	Stop() error
	IsRunning() bool
}

// State representa o estado atual do worker pool.
type State int

const (
	StartIdle    State = iota // Estado inicial, ocioso.
	StateRunning              // Estado em execução.
	StateStopped              // Estado parado.
)

// Config contém a configuração do worker pool.
type Config struct {
	WorkerCount int          // Número de workers.
	Logger      *slog.Logger // Logger para registrar eventos.
}

// DefaultConfig retorna a configuração padrão do worker pool.
func DefaultConfig() Config {
	return Config{
		WorkerCount: 1,
		Logger:      slog.Default(),
	}
}

// workerPool implementa a interface WorkerPool.
type workerPool struct {
	WorkerCount int
	ProcessFunc ProcessFunc
	logger      *slog.Logger
	state       State
	stateMutex  sync.Mutex
	stopCh      chan struct{}
	stopWg      sync.WaitGroup
}

// New cria um novo worker pool com a função de processamento e configuração fornecidas.
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

// Start inicia o pool de workers e retorna um canal de resultados.
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

	// Cria os workers e os inicia em goroutines separadas.
	for i := 0; i < wp.WorkerCount; i++ {
		go wp.worker(ctx, i, inputCh, resultCh)
	}

	// Goroutine para aguardar a finalização dos workers e fechar o canal de resultados.
	go func() {
		wp.stopWg.Wait()
		close(resultCh)

		wp.stateMutex.Lock()
		wp.state = StartIdle
		wp.stateMutex.Unlock()
	}()

	return resultCh, nil
}

// Stop finaliza o pool de workers de forma controlada.
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

// IsRunning verifica se o worker pool está em execução.
func (wp *workerPool) IsRunning() bool {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()
	return wp.state == StateRunning
}

// worker representa um trabalhador individual que processa trabalhos do canal de entrada.
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

			// Processa o trabalho recebido.
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
