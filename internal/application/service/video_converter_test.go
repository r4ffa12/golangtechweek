package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFFmpegService é um mock para o serviço FFmpeg
type MockFFmpegService struct {
	mock.Mock
}

func (m *MockFFmpegService) ConvertToHLS(ctx context.Context, input string, outputDir string) ([]OutputFile, error) {
	args := m.Called(ctx, input, outputDir)
	return args.Get(0).([]OutputFile), args.Error(1)
}

// MockVideoRepository é um mock para o repositório de vídeos
type MockVideoRepository struct {
	mock.Mock
}

func (m *MockVideoRepository) Create(ctx context.Context, video *entity.Video) error {
	args := m.Called(ctx, video)
	return args.Error(0)
}

func (m *MockVideoRepository) FindByID(ctx context.Context, id string) (*entity.Video, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) List(ctx context.Context, page, pageSize int) ([]*entity.Video, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*entity.Video), args.Error(1)
}

func (m *MockVideoRepository) UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error {
	args := m.Called(ctx, id, status, errorMessage)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateHLSPath(ctx context.Context, id string, hlsPath, manifestPath string) error {
	args := m.Called(ctx, id, hlsPath, manifestPath)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateS3Status(ctx context.Context, id string, uploadStatus string) error {
	args := m.Called(ctx, id, uploadStatus)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateS3URLs(ctx context.Context, id string, s3URL, s3ManifestURL string) error {
	args := m.Called(ctx, id, s3URL, s3ManifestURL)
	return args.Error(0)
}

func (m *MockVideoRepository) UpdateS3Keys(ctx context.Context, id string, segmentKey string, manifestKey string) error {
	args := m.Called(ctx, id, segmentKey, manifestKey)
	return args.Error(0)
}

func (m *MockVideoRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Substituir o FFmpegService no VideoConverterService para testes
func replaceFFmpegService(converter *VideoConverterService, ffmpeg FFmpegServiceInterface) {
	converter.ffmpeg = ffmpeg
}

func TestNewVideoConverter(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := VideoConverterConfig{
		WorkerCount: 3,
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})),
	}

	// Act
	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Assert
	assert.NotNil(t, converter)
	assert.Equal(t, mockRepo, converter.videoRepo)
	assert.Equal(t, mockFFmpeg, converter.ffmpeg)
	assert.NotNil(t, converter.workerPool)
	assert.NotNil(t, converter.logger)
}

func TestVideoConverterService_StartConversion_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()
	config.WorkerCount = 1 // Usar apenas 1 worker para simplificar o teste

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Configurar o mock do repositório para retornar sucesso ao atualizar o status
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusProcessing, "").Return(nil)
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusCompleted, "").Return(nil)
	mockRepo.On("UpdateHLSPath", mock.Anything, "test-video-id", "output/dir", "output/dir/manifest.m3u8").Return(nil)

	// Configurar o mock do FFmpeg para retornar arquivos de saída simulados
	outputFiles := []OutputFile{
		{Path: "output/dir/manifest.m3u8", Type: entity.FileTypeManifest},
		{Path: "output/dir/segment_0.ts", Type: entity.FileTypeSegment},
	}
	mockFFmpeg.On("ConvertToHLS", mock.Anything, "input/path", mock.Anything).Return(outputFiles, nil)

	// Criar um canal de entrada com capacidade para evitar bloqueio
	inputCh := make(chan ConversionJob, 1)

	// Criar um contexto com timeout para evitar que o teste fique preso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act
	resultCh, err := converter.StartConversion(ctx, inputCh)

	// Assert - verificar se o serviço iniciou corretamente
	assert.NoError(t, err)
	assert.NotNil(t, resultCh)

	// Enviar um job para o canal de entrada
	inputCh <- ConversionJob{
		VideoID:   "test-video-id",
		InputPath: "input/path",
		OutputDir: "output/dir",
	}

	// Fechar o canal de entrada para sinalizar que não há mais jobs
	close(inputCh)

	// Ler o resultado do canal de saída
	result := <-resultCh

	// Verificar o resultado
	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
	assert.Equal(t, "test-video-id", result.VideoID)
	assert.NotEmpty(t, result.OutputFiles)
	assert.Greater(t, result.Duration, time.Duration(0))

	// Verificar se os mocks foram chamados conforme esperado
	mockRepo.AssertExpectations(t)
	mockFFmpeg.AssertExpectations(t)

	// Parar o serviço apenas se ainda estiver em execução
	if converter.IsRunning() {
		err = converter.StopConversion()
		assert.NoError(t, err)
	}
}

func TestVideoConverterService_StartConversion_FFmpegError(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()
	config.WorkerCount = 1 // Usar apenas 1 worker para simplificar o teste

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Configurar o mock do repositório
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusProcessing, "").Return(nil)
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusError, mock.Anything).Return(nil)

	// Configurar o mock do FFmpeg para retornar erro
	ffmpegError := errors.New("erro na conversão")
	mockFFmpeg.On("ConvertToHLS", mock.Anything, "input/path", mock.Anything).Return([]OutputFile{}, ffmpegError)

	// Criar um canal de entrada com capacidade para evitar bloqueio
	inputCh := make(chan ConversionJob, 1)

	// Criar um contexto com timeout para evitar que o teste fique preso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act
	resultCh, err := converter.StartConversion(ctx, inputCh)

	// Assert - verificar se o serviço iniciou corretamente
	assert.NoError(t, err)
	assert.NotNil(t, resultCh)

	// Enviar um job para o canal de entrada
	inputCh <- ConversionJob{
		VideoID:   "test-video-id",
		InputPath: "input/path",
		OutputDir: "output/dir",
	}

	// Fechar o canal de entrada para sinalizar que não há mais jobs
	close(inputCh)

	// Ler o resultado do canal de saída
	result := <-resultCh

	// Verificar o resultado
	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "test-video-id", result.VideoID)
	assert.Contains(t, result.Error.Error(), "erro na conversão")

	// Verificar se os mocks foram chamados conforme esperado
	mockRepo.AssertExpectations(t)
	mockFFmpeg.AssertExpectations(t)

	// Parar o serviço apenas se ainda estiver em execução
	if converter.IsRunning() {
		err = converter.StopConversion()
		assert.NoError(t, err)
	}
}

func TestVideoConverterService_StartConversion_UpdateStatusError(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()
	config.WorkerCount = 1 // Usar apenas 1 worker para simplificar o teste

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Configurar o mock do repositório para retornar erro ao atualizar o status
	updateError := errors.New("erro ao atualizar status")
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusProcessing, "").Return(updateError)
	mockRepo.On("UpdateStatus", mock.Anything, "test-video-id", entity.StatusError, mock.Anything).Return(nil)

	// Criar um canal de entrada com capacidade para evitar bloqueio
	inputCh := make(chan ConversionJob, 1)

	// Criar um contexto com timeout para evitar que o teste fique preso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act
	resultCh, err := converter.StartConversion(ctx, inputCh)

	// Assert - verificar se o serviço iniciou corretamente
	assert.NoError(t, err)
	assert.NotNil(t, resultCh)

	// Enviar um job para o canal de entrada
	inputCh <- ConversionJob{
		VideoID:   "test-video-id",
		InputPath: "input/path",
		OutputDir: "output/dir",
	}

	// Fechar o canal de entrada para sinalizar que não há mais jobs
	close(inputCh)

	// Ler o resultado do canal de saída
	result := <-resultCh

	// Verificar o resultado
	assert.False(t, result.Success)
	assert.NotNil(t, result.Error)
	assert.Equal(t, "test-video-id", result.VideoID)
	assert.Contains(t, result.Error.Error(), "erro ao atualizar status")

	// Verificar se os mocks foram chamados conforme esperado
	mockRepo.AssertExpectations(t)

	// Parar o serviço apenas se ainda estiver em execução
	if converter.IsRunning() {
		err = converter.StopConversion()
		assert.NoError(t, err)
	}
}

func TestVideoConverterService_StartConversion_AlreadyRunning(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Criar canais de entrada
	inputCh1 := make(chan ConversionJob, 1)
	inputCh2 := make(chan ConversionJob, 1)

	// Criar um contexto com timeout para evitar que o teste fique preso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act - Iniciar o serviço pela primeira vez
	resultCh1, err1 := converter.StartConversion(ctx, inputCh1)

	// Assert - verificar se o serviço iniciou corretamente
	assert.NoError(t, err1)
	assert.NotNil(t, resultCh1)

	// Act - Tentar iniciar o serviço novamente
	resultCh2, err2 := converter.StartConversion(ctx, inputCh2)

	// Assert - verificar se o serviço retornou erro ao tentar iniciar novamente
	assert.Error(t, err2)
	assert.Nil(t, resultCh2)
	assert.Contains(t, err2.Error(), "o serviço de conversão já está em execução")

	// Parar o serviço
	err := converter.StopConversion()
	assert.NoError(t, err)

	// Fechar os canais de entrada
	close(inputCh1)
	close(inputCh2)
}

func TestVideoConverterService_StopConversion_NotRunning(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Act - Tentar parar o serviço que não está em execução
	err := converter.StopConversion()

	// Assert - verificar se o serviço retornou erro ao tentar parar
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "o serviço de conversão não está em execução")
}

func TestVideoConverterService_IsRunning(t *testing.T) {
	// Arrange
	mockRepo := new(MockVideoRepository)
	mockFFmpeg := new(MockFFmpegService)
	config := DefaultVideoConverterConfig()

	converter := NewVideoConverter(mockFFmpeg, mockRepo, config)

	// Assert - verificar se o serviço não está em execução inicialmente
	assert.False(t, converter.IsRunning())

	// Criar um canal de entrada
	inputCh := make(chan ConversionJob, 1)

	// Criar um contexto com timeout para evitar que o teste fique preso
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Act - Iniciar o serviço
	resultCh, err := converter.StartConversion(ctx, inputCh)

	// Assert - verificar se o serviço iniciou corretamente
	assert.NoError(t, err)
	assert.NotNil(t, resultCh)

	// Assert - verificar se o serviço está em execução após iniciar
	assert.True(t, converter.IsRunning())

	// Act - Parar o serviço
	err = converter.StopConversion()

	// Assert - verificar se o serviço parou corretamente
	assert.NoError(t, err)

	// Assert - verificar se o serviço não está mais em execução após parar
	assert.False(t, converter.IsRunning())

	// Fechar o canal de entrada
	close(inputCh)
}
