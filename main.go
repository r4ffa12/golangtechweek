package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/r4ffa12/golangtechweek/pkg/wokerpool"
)

type NumeroJob struct {
	Numero int
}
type ResultadoNumero struct {
	Valor     int
	WorkerID  int
	Timestamp time.Time
}

func processarNumero(ctx context.Context, job wokerpool.Job) wokerpool.Result {
	numero := job.(NumeroJob).Numero
	workerID := numero % 3

	sleepTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(sleepTime)

	return ResultadoNumero{
		Valor:     numero,
		WorkerID:  workerID,
		Timestamp: time.Now(),
	}
}

func main() {
	valorMaximo := 20
	bufferSize := 10

	pool := wokerpool.New(processarNumero, wokerpool.Config{
		WorkerCount: 3,
	})

	inputCh := make(chan wokerpool.Job, bufferSize)
	ctx := context.Background()

	resultCh, err := pool.Start(ctx, inputCh)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(valorMaximo)

	fmt.Println("Iniciando o pool de workers", valorMaximo, "numeros")

	go func() {
		for i := 0; i < valorMaximo; i++ {
			inputCh <- NumeroJob{Numero: i}
		}
		close(inputCh)
	}()

	go func() {
		for result := range resultCh {
			r := result.(ResultadoNumero)
			fmt.Printf("Numero: %d, WorkerID: %d, Timestamp: %s\n", r.Valor, r.WorkerID, r.Timestamp.Format(time.RFC3339))
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Println(" \n Pool de workers finalizadon\n Todos os numeros foram processados\n", valorMaximo)

}
