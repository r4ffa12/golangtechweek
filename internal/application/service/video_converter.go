package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	"github.com/devfullcycle/golangtechweek/internal/domain/repository"
	"github.com/devfullcycle/golangtechweek/pkg/workerpool"
)

// ConversionJob representa um trabalho de conversão de vídeo
type ConversionJob struct {
	VideoID   string // ID do vídeo no banco de dados
	InputPath string // Caminho do arquivo de entrada
	OutputDir string // Diretório de saída para os arquivos convertidos
}

// ConversionResult representa o resultado de uma conversão
type ConversionResult struct {
	VideoID     string        // ID do vídeo no banco de dados
	Success     bool          // Indica se a conversão foi bem-sucedida
	Error       error         // Erro ocorrido durante a conversão, se houver
	OutputFiles []OutputFile  // Lista de arquivos gerados pela conversão
	Duration    time.Duration // Duração do processo de conversão
}

// VideoConverterService implementa o serviço de conversão de vídeos
type VideoConverterService struct {
	ffmpeg     FFmpegServiceInterface
	videoRepo  repository.VideoRepository
	workerPool workerpool.WorkerPool
	logger     *slog.Logger
}

// VideoConverterConfig representa a configuração do serviço de conversão de vídeo
type VideoConverterConfig struct {
	WorkerCount int          // Número de workers para processamento paralelo
	Logger      *slog.Logger // Logger para registro de eventos
}

// DefaultVideoConverterConfig retorna uma configuração padrão para o serviço
func DefaultVideoConverterConfig() VideoConverterConfig {
	return VideoConverterConfig{
		WorkerCount: 3,
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}
}

// NewVideoConverter cria uma nova instância do serviço de conversão de vídeos
func NewVideoConverter(ffmpeg FFmpegServiceInterface, videoRepo repository.VideoRepository, config VideoConverterConfig) *VideoConverterService {
	if config.WorkerCount <= 0 {
		config.WorkerCount = 1
	}

	if config.Logger == nil {
		config.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	service := &VideoConverterService{
		ffmpeg:    ffmpeg,
		videoRepo: videoRepo,
		logger:    config.Logger,
	}

	// Cria a função de processamento para o worker pool
	processFunc := func(ctx context.Context, job workerpool.Job) workerpool.Result {
		// Conversão de tipo segura
		conversionJob, ok := job.(ConversionJob)
		if !ok {
			return ConversionResult{
				Success: false,
				Error:   fmt.Errorf("job inválido: esperado ConversionJob"),
			}
		}

		return service.processJob(ctx, conversionJob)
	}

	// Cria a configuração para o worker pool
	wpConfig := workerpool.Config{
		WorkerCount: config.WorkerCount,
		Logger:      config.Logger,
	}

	// Cria o worker pool
	service.workerPool = workerpool.New(processFunc, wpConfig)

	return service
}

// StartConversion inicia o serviço de conversão de vídeos
func (c *VideoConverterService) StartConversion(ctx context.Context, inputCh <-chan ConversionJob) (<-chan ConversionResult, error) {
	// Verifica se o serviço já está em execução
	if c.workerPool.IsRunning() {
		return nil, fmt.Errorf("o serviço de conversão já está em execução")
	}

	// Cria um canal para adaptar o canal de entrada genérico para o tipo específico
	jobCh := make(chan workerpool.Job)

	// Goroutine para adaptar o canal de entrada
	go func() {
		defer close(jobCh)
		for job := range inputCh {
			select {
			case jobCh <- job:
				// Job enviado com sucesso
			case <-ctx.Done():
				return
			}
		}
	}()

	// Inicia o worker pool
	resultCh, err := c.workerPool.Start(ctx, jobCh)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar o worker pool: %w", err)
	}

	// Cria um canal para adaptar o canal de saída genérico para o tipo específico
	conversionResultCh := make(chan ConversionResult)

	// Goroutine para adaptar o canal de saída
	go func() {
		defer close(conversionResultCh)
		for result := range resultCh {
			// Conversão de tipo segura
			convResult, ok := result.(ConversionResult)
			if !ok {
				c.logger.Error("resultado inválido do worker pool", "error", "tipo incompatível")
				continue
			}

			select {
			case conversionResultCh <- convResult:
				// Resultado enviado com sucesso
			case <-ctx.Done():
				return
			}
		}
	}()

	return conversionResultCh, nil
}

// StopConversion interrompe o serviço de conversão
func (c *VideoConverterService) StopConversion() error {
	// Verifica se o serviço está em execução
	if !c.workerPool.IsRunning() {
		return fmt.Errorf("o serviço de conversão não está em execução")
	}

	// Para o worker pool
	return c.workerPool.Stop()
}

// IsRunning retorna true se o serviço estiver em execução
func (c *VideoConverterService) IsRunning() bool {
	return c.workerPool.IsRunning()
}

// processJob processa um trabalho de conversão de vídeo
func (c *VideoConverterService) processJob(ctx context.Context, job ConversionJob) ConversionResult {
	startTime := time.Now()
	c.logger.Info("Iniciando processamento de vídeo", "video_id", job.VideoID)

	// Inicializa o resultado com falha por padrão
	result := ConversionResult{
		VideoID:  job.VideoID,
		Success:  false,
		Duration: time.Since(startTime),
	}

	// Etapa 1: Atualiza o status do vídeo para "processing"
	if err := c.updateVideoStatusToProcessing(ctx, job.VideoID); err != nil {
		result.Error = err
		return result
	}

	// Etapa 2: Prepara o diretório de saída
	outputDir := c.prepareOutputDirectory(job)

	// Etapa 3: Converte o vídeo para HLS
	outputFiles, err := c.convertVideoToHLS(ctx, job.VideoID, job.InputPath, outputDir)
	if err != nil {
		result.Error = err
		return result
	}

	// Etapa 4: Atualiza o resultado com sucesso
	result.Success = true
	result.OutputFiles = outputFiles
	result.Duration = time.Since(startTime)

	// Etapa 5: Processa os arquivos de saída e atualiza o banco de dados
	c.processOutputFiles(ctx, job.VideoID, outputFiles)

	c.logger.Info("Processamento de vídeo concluído com sucesso",
		"video_id", job.VideoID,
		"duration", result.Duration.String(),
		"output_files", len(result.OutputFiles))

	return result
}

// updateVideoStatusToProcessing atualiza o status do vídeo para "processing"
func (c *VideoConverterService) updateVideoStatusToProcessing(ctx context.Context, videoID string) error {
	err := c.videoRepo.UpdateStatus(ctx, videoID, entity.StatusProcessing, "")
	if err != nil {
		errWithContext := fmt.Errorf("erro ao atualizar status do vídeo para processing: %w", err)
		c.logger.Error("Erro ao atualizar status do vídeo", "video_id", videoID, "error", err)
		c.videoRepo.UpdateStatus(ctx, videoID, entity.StatusError, errWithContext.Error())
		return errWithContext
	}
	return nil
}

// prepareOutputDirectory prepara o diretório de saída para os arquivos convertidos
func (c *VideoConverterService) prepareOutputDirectory(job ConversionJob) string {
	outputDir := job.OutputDir
	if outputDir == "" {
		outputDir = filepath.Join("uploads", "converted", job.VideoID)
	}
	return outputDir
}

// convertVideoToHLS converte o vídeo para o formato HLS
func (c *VideoConverterService) convertVideoToHLS(ctx context.Context, videoID, inputPath, outputDir string) ([]OutputFile, error) {
	outputFiles, err := c.ffmpeg.ConvertToHLS(ctx, inputPath, outputDir)
	if err != nil {
		errWithContext := fmt.Errorf("erro ao converter vídeo para HLS: %w", err)
		c.logger.Error("Erro ao converter vídeo para HLS", "video_id", videoID, "error", err)
		c.videoRepo.UpdateStatus(ctx, videoID, entity.StatusError, errWithContext.Error())
		return nil, errWithContext
	}
	return outputFiles, nil
}

// processOutputFiles processa os arquivos de saída e atualiza o banco de dados
func (c *VideoConverterService) processOutputFiles(ctx context.Context, videoID string, outputFiles []OutputFile) {
	// Encontra o manifesto e os segmentos
	manifestPath, hlsPath := c.findManifestAndHLSPaths(outputFiles)

	// Atualiza os caminhos HLS e Manifest no banco de dados
	if manifestPath != "" && hlsPath != "" {
		c.updateHLSPaths(ctx, videoID, hlsPath, manifestPath)
	}

	// Atualiza o status do vídeo para "completed"
	c.updateVideoStatusToCompleted(ctx, videoID)
}

// findManifestAndHLSPaths encontra os caminhos do manifesto e do diretório HLS
// Retorna o caminho do manifesto e o caminho do diretório HLS
func (c *VideoConverterService) findManifestAndHLSPaths(outputFiles []OutputFile) (string, string) {
	var manifestPath string
	var hlsPath string

	// Itera sobre os arquivos até encontrar tanto o manifesto quanto o primeiro segmento
	for _, file := range outputFiles {
		// Se ainda não encontramos o manifesto e este arquivo é um manifesto
		if manifestPath == "" && file.Type == entity.FileTypeManifest {
			manifestPath = file.Path
		}

		// Se ainda não encontramos o caminho HLS e este arquivo é um segmento
		if hlsPath == "" && file.Type == entity.FileTypeSegment {
			// Usa o diretório do primeiro segmento como caminho HLS
			hlsPath = filepath.Dir(file.Path)
		}

		// Se já encontramos ambos, podemos sair do loop
		if manifestPath != "" && hlsPath != "" {
			break
		}
	}

	return manifestPath, hlsPath
}

// updateHLSPaths atualiza os caminhos HLS e Manifest no banco de dados
func (c *VideoConverterService) updateHLSPaths(ctx context.Context, videoID, hlsPath, manifestPath string) {
	err := c.videoRepo.UpdateHLSPath(ctx, videoID, hlsPath, manifestPath)
	if err != nil {
		c.logger.Error("Erro ao atualizar caminhos HLS", "video_id", videoID, "error", err)
		// Não falha a conversão por erro na atualização dos caminhos
	}
}

// updateVideoStatusToCompleted atualiza o status do vídeo para "completed"
func (c *VideoConverterService) updateVideoStatusToCompleted(ctx context.Context, videoID string) {
	err := c.videoRepo.UpdateStatus(ctx, videoID, entity.StatusCompleted, "")
	if err != nil {
		c.logger.Error("Erro ao atualizar status do vídeo para completed", "video_id", videoID, "error", err)
		// Não falha a conversão por erro na atualização do status
	}
}
