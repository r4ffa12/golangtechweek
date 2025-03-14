//go:build integration
// +build integration

// Este arquivo contém testes de integração que acessam o banco de dados real.
// Para executar estes testes, use o comando:
// go test -tags=integration ./internal/infra/database/repository

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	"github.com/devfullcycle/golangtechweek/internal/infra/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type VideoRepositoryTestSuite struct {
	suite.Suite
	db         *sql.DB
	repository *VideoRepositoryPostgres
	ctx        context.Context
}

func (suite *VideoRepositoryTestSuite) SetupSuite() {
	// Configuração do banco de dados de teste
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "conversorgo"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	var err error
	suite.db, err = database.NewConnection(dbConfig)
	if err != nil {
		suite.T().Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Limpar a tabela de vídeos antes de iniciar os testes
	_, err = suite.db.Exec("DELETE FROM videos")
	if err != nil {
		suite.T().Fatalf("Erro ao limpar a tabela de vídeos: %v", err)
	}

	suite.repository = NewVideoRepositoryPostgres(suite.db)
	suite.ctx = context.Background()
}

func (suite *VideoRepositoryTestSuite) TearDownSuite() {
	// Limpar a tabela de vídeos após os testes
	_, err := suite.db.Exec("DELETE FROM videos")
	if err != nil {
		suite.T().Fatalf("Erro ao limpar a tabela de vídeos: %v", err)
	}

	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *VideoRepositoryTestSuite) TestCreate() {
	video := entity.NewVideo("Teste de Vídeo", "/path/to/video.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Verificar se o vídeo foi criado corretamente
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM videos WHERE id = $1", video.ID).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)
}

func (suite *VideoRepositoryTestSuite) TestFindByID() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Busca", "/path/to/search.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Buscar o vídeo pelo ID
	foundVideo, err := suite.repository.FindByID(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundVideo)
	assert.Equal(suite.T(), video.ID, foundVideo.ID)
	assert.Equal(suite.T(), video.Title, foundVideo.Title)
	assert.Equal(suite.T(), video.FilePath, foundVideo.FilePath)
	assert.Equal(suite.T(), video.Status, foundVideo.Status)
}

func (suite *VideoRepositoryTestSuite) TestFindByIDNotFound() {
	// Buscar um vídeo com ID inexistente, mas com formato UUID válido
	_, err := suite.repository.FindByID(suite.ctx, uuid.New().String())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrVideoNotFound, err)
}

func (suite *VideoRepositoryTestSuite) TestList() {
	// Limpar a tabela para garantir um estado conhecido
	_, err := suite.db.Exec("DELETE FROM videos")
	assert.NoError(suite.T(), err)

	// Criar vários vídeos para o teste
	for i := 1; i <= 15; i++ {
		video := entity.NewVideo(fmt.Sprintf("Vídeo %d", i), fmt.Sprintf("/path/to/video%d.mp4", i))
		err := suite.repository.Create(suite.ctx, video)
		assert.NoError(suite.T(), err)
	}

	// Testar a primeira página
	videos, err := suite.repository.List(suite.ctx, 1, 10)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), videos, 10)

	// Testar a segunda página
	videos, err = suite.repository.List(suite.ctx, 2, 10)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), videos, 5)
}

func (suite *VideoRepositoryTestSuite) TestUpdateStatus() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Atualização de Status", "/path/to/status.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Atualizar o status
	err = suite.repository.UpdateStatus(suite.ctx, video.ID, entity.StatusProcessing, "")
	assert.NoError(suite.T(), err)

	// Verificar se o status foi atualizado
	foundVideo, err := suite.repository.FindByID(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entity.StatusProcessing, foundVideo.Status)
}

func (suite *VideoRepositoryTestSuite) TestUpdateHLSPath() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Atualização de HLS", "/path/to/hls.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Atualizar os caminhos HLS
	hlsPath := "/path/to/hls"
	manifestPath := "/path/to/manifest.m3u8"
	err = suite.repository.UpdateHLSPath(suite.ctx, video.ID, hlsPath, manifestPath)
	assert.NoError(suite.T(), err)

	// Verificar se os caminhos foram atualizados
	foundVideo, err := suite.repository.FindByID(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), hlsPath, foundVideo.HLSPath)
	assert.Equal(suite.T(), manifestPath, foundVideo.ManifestPath)
}

func (suite *VideoRepositoryTestSuite) TestUpdateS3Status() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Atualização de S3 Status", "/path/to/s3status.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Atualizar o status de upload
	err = suite.repository.UpdateS3Status(suite.ctx, video.ID, entity.UploadStatusPendingS3)
	assert.NoError(suite.T(), err)

	// Verificar se o status foi atualizado
	foundVideo, err := suite.repository.FindByID(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entity.UploadStatusPendingS3, foundVideo.UploadStatus)
}

func (suite *VideoRepositoryTestSuite) TestUpdateS3URLs() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Atualização de S3 URLs", "/path/to/s3urls.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Atualizar as URLs do S3
	s3URL := "https://bucket.s3.amazonaws.com/videos/123"
	s3ManifestURL := "https://bucket.s3.amazonaws.com/manifests/123.m3u8"
	err = suite.repository.UpdateS3URLs(suite.ctx, video.ID, s3URL, s3ManifestURL)
	assert.NoError(suite.T(), err)

	// Verificar se as URLs foram atualizadas
	foundVideo, err := suite.repository.FindByID(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), s3URL, foundVideo.S3URL)
	assert.Equal(suite.T(), s3ManifestURL, foundVideo.S3ManifestURL)
}

func (suite *VideoRepositoryTestSuite) TestUpdateS3Keys() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Atualização de S3 Keys", "/path/to/s3keys.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Atualizar as chaves do S3
	segmentKey := "videos/123"
	manifestKey := "manifests/123.m3u8"
	err = suite.repository.UpdateS3Keys(suite.ctx, video.ID, segmentKey, manifestKey)
	assert.NoError(suite.T(), err)

	// Verificar se as chaves foram atualizadas
	var foundSegmentKey, foundManifestKey string
	err = suite.db.QueryRow("SELECT segment_key, manifest_key FROM videos WHERE id = $1", video.ID).Scan(&foundSegmentKey, &foundManifestKey)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), segmentKey, foundSegmentKey)
	assert.Equal(suite.T(), manifestKey, foundManifestKey)
}

func (suite *VideoRepositoryTestSuite) TestDelete() {
	// Criar um vídeo para o teste
	video := entity.NewVideo("Teste de Exclusão", "/path/to/delete.mp4")
	err := suite.repository.Create(suite.ctx, video)
	assert.NoError(suite.T(), err)

	// Excluir o vídeo
	err = suite.repository.Delete(suite.ctx, video.ID)
	assert.NoError(suite.T(), err)

	// Verificar se o vídeo foi marcado como excluído
	var deletedAt sql.NullTime
	err = suite.db.QueryRow("SELECT deleted_at FROM videos WHERE id = $1", video.ID).Scan(&deletedAt)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), deletedAt.Valid)

	// Tentar buscar o vídeo excluído
	_, err = suite.repository.FindByID(suite.ctx, video.ID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), ErrVideoNotFound, err)
}

func TestVideoRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(VideoRepositoryTestSuite))
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão se não estiver definida
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
