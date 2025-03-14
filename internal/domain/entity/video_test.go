package entity

import (
	"testing"
	"time"
)

func NewVideTest(t *testing.T) {
	title := "Meu Vídeo de Teste"
	filePath := "/tmp/video.mp4"

	video := NewVideo(title, filePath)

	// Verifica se o ID foi gerado
	if video.ID == "" {
		t.Error("ID não deveria ser vazio")
	}

	// Verifica se os campos foram preenchidos corretamente
	if video.Title != title {
		t.Errorf("Esperado Title %s, obtido %s", title, video.Title)
	}

	if video.FilePath != filePath {
		t.Errorf("Esperado FilePath %s, obtido %s", filePath, video.FilePath)
	}

	// Verifica se o status inicial está correto
	if video.Status != StatusPending {
		t.Errorf("Esperado Status %s, obtido %s", StatusPending, video.Status)
	}

	// Verifica se o status de upload inicial está correto
	if video.UploadStatus != UploadStatusNone {
		t.Errorf("Esperado UploadStatus %s, obtido %s", UploadStatusNone, video.UploadStatus)
	}

	// Verifica se as datas foram preenchidas
	if video.CreatedAt.IsZero() {
		t.Error("CreatedAt não deveria ser zero")
	}

	if video.UpdatedAt.IsZero() {
		t.Error("UpdatedAt não deveria ser zero")
	}

	// Verifica se os campos opcionais estão vazios
	if video.HLSPath != "" {
		t.Errorf("Esperado HLSPath vazio, obtido %s", video.HLSPath)
	}

	if video.ManifestPath != "" {
		t.Errorf("Esperado ManifestPath vazio, obtido %s", video.ManifestPath)
	}

	if video.S3URL != "" {
		t.Errorf("Esperado S3URL vazio, obtido %s", video.S3URL)
	}

	if video.S3ManifestURL != "" {
		t.Errorf("Esperado S3ManifestURL vazio, obtido %s", video.S3ManifestURL)
	}

	if video.ErrorMessage != "" {
		t.Errorf("Esperado ErrorMessage vazio, obtido %s", video.ErrorMessage)
	}
}

func TestMarkAsProcessing(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	oldUpdatedAt := video.UpdatedAt

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	video.MarkAsProcessing()

	if video.Status != StatusProcessing {
		t.Errorf("Esperado Status %s, obtido %s", StatusProcessing, video.Status)
	}

	if !video.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt deveria ter sido atualizado")
	}
}

func TestMarkAsCompleted(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	oldUpdatedAt := video.UpdatedAt

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	hlsPath := "/tmp/output/123"
	manifestPath := "/tmp/output/123/playlist.m3u8"

	video.MarkAsCompleted(hlsPath, manifestPath)

	if video.Status != StatusCompleted {
		t.Errorf("Esperado Status %s, obtido %s", StatusCompleted, video.Status)
	}

	if video.HLSPath != hlsPath {
		t.Errorf("Esperado HLSPath %s, obtido %s", hlsPath, video.HLSPath)
	}

	if video.ManifestPath != manifestPath {
		t.Errorf("Esperado ManifestPath %s, obtido %s", manifestPath, video.ManifestPath)
	}

	if !video.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt deveria ter sido atualizado")
	}
}

func TestMarkAsFailed(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	oldUpdatedAt := video.UpdatedAt
	errorMsg := "Erro ao processar vídeo"

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	video.MarkAsFailed(errorMsg)

	if video.Status != StatusError {
		t.Errorf("Esperado Status %s, obtido %s", StatusError, video.Status)
	}

	if video.ErrorMessage != errorMsg {
		t.Errorf("Esperado ErrorMessage %s, obtido %s", errorMsg, video.ErrorMessage)
	}

	if !video.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt deveria ter sido atualizado")
	}
}

func TestSetS3URL(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	oldUpdatedAt := video.UpdatedAt
	url := "https://bucket.s3.amazonaws.com/videos/123/video.m3u8"

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	video.SetS3URL(url)

	if video.S3URL != url {
		t.Errorf("Esperado S3URL %s, obtido %s", url, video.S3URL)
	}

	if !video.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt deveria ter sido atualizado")
	}
}

func TestSetS3ManifestURL(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	oldUpdatedAt := video.UpdatedAt
	url := "https://bucket.s3.amazonaws.com/videos/123/playlist.m3u8"

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	video.SetS3ManifestURL(url)

	if video.S3ManifestURL != url {
		t.Errorf("Esperado S3ManifestURL %s, obtido %s", url, video.S3ManifestURL)
	}

	if !video.UpdatedAt.After(oldUpdatedAt) {
		t.Error("UpdatedAt deveria ter sido atualizado")
	}
}

func TestIsCompleted(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")

	if video.IsCompleted() {
		t.Error("Vídeo não deveria estar completo inicialmente")
	}

	hlsPath := "/tmp/output/123"
	manifestPath := "/tmp/output/123/playlist.m3u8"

	video.MarkAsCompleted(hlsPath, manifestPath)

	if !video.IsCompleted() {
		t.Error("Vídeo deveria estar completo após MarkAsCompleted")
	}
}

func TestGetHLSDirectory(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	hlsPath := "/tmp/output/123"
	video.HLSPath = hlsPath

	if video.GetHLSDirectory() != hlsPath {
		t.Errorf("Esperado HLSPath %s, obtido %s", hlsPath, video.GetHLSDirectory())
	}
}

func TestGetManifestPath(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	manifestPath := "/tmp/output/123/playlist.m3u8"
	video.ManifestPath = manifestPath

	if video.GetManifestPath() != manifestPath {
		t.Errorf("Esperado ManifestPath %s, obtido %s", manifestPath, video.GetManifestPath())
	}
}

func TestGenerateOutputPath(t *testing.T) {
	video := NewVideo("Test Video", "/tmp/video.mp4")
	baseDir := "/tmp/output"

	oldUpdatedAt := video.UpdatedAt

	// Aguarda um momento para garantir que o timestamp seja diferente
	time.Sleep(1 * time.Millisecond)

	outputPath := video.GenerateOutputPath(baseDir)

	expectedOutputPath := baseDir + "/converted/" + video.ID
	if outputPath != expectedOutputPath {
		t.Errorf("Esperado outputPath %s, obtido %s", expectedOutputPath, outputPath)
	}

	// Não deve mais atualizar os campos da entidade
	if video.HLSPath != "" {
		t.Errorf("HLSPath não deveria ter sido atualizado, obtido %s", video.HLSPath)
	}

	if video.ManifestPath != "" {
		t.Errorf("ManifestPath não deveria ter sido atualizado, obtido %s", video.ManifestPath)
	}

	// Não deve mais atualizar o timestamp
	if video.UpdatedAt != oldUpdatedAt {
		t.Error("UpdatedAt não deveria ter sido atualizado")
	}
}
