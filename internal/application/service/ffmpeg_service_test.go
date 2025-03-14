//go:build integration
// +build integration

package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/devfullcycle/golangtechweek/internal/application/service"
	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFFmpegService_ConvertToHLS_Integration é um teste de integração
// que requer a presença de um arquivo de vídeo real e o FFmpeg instalado
func TestFFmpegService_ConvertToHLS_Integration(t *testing.T) {
	// Criar instância do serviço
	ffmpegService := service.NewFFmpegService()

	// Definir diretório de saída
	outputDir := "/app/uploads/converted"

	// Limpar o diretório de saída antes do teste
	os.RemoveAll(outputDir)

	// Criar o diretório de saída
	err := os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	// Garantir que o diretório será limpo após o teste
	defer os.RemoveAll(outputDir)

	// Caminho para o vídeo de teste na pasta uploads
	testVideoPath := "/app/uploads/44444444-4444-4444-4444-444444444444.mp4"

	// Verificar se o arquivo de teste existe
	_, err = os.Stat(testVideoPath)
	if os.IsNotExist(err) {
		t.Logf("Arquivo de teste não encontrado: %s", testVideoPath)
		t.Skip("Arquivo de vídeo de teste não encontrado")
	}

	// Executar a conversão
	ctx := context.Background()
	outputFiles, err := ffmpegService.ConvertToHLS(ctx, testVideoPath, outputDir)

	// Verificar se não houve erro
	require.NoError(t, err)

	// Verificar se foram gerados arquivos
	assert.NotEmpty(t, outputFiles)

	// Verificar se o manifesto foi gerado
	var hasManifest bool
	for _, file := range outputFiles {
		if file.Type == entity.FileTypeManifest {
			hasManifest = true
			break
		}
	}
	assert.True(t, hasManifest, "O manifesto não foi gerado")

	// Verificar se foram gerados segmentos
	var segmentCount int
	for _, file := range outputFiles {
		if file.Type == entity.FileTypeSegment {
			segmentCount++
		}
	}
	assert.Greater(t, segmentCount, 0, "Nenhum segmento foi gerado")

	// Imprimir informações sobre os arquivos gerados
	t.Logf("Arquivos gerados no diretório: %s", outputDir)
	for _, file := range outputFiles {
		t.Logf("- %s (tipo: %s)", file.Path, file.Type)
	}
}
