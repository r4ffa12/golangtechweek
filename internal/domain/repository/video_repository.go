package repository

import (
	"context"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
)

// VideoRepository define as operações que podem ser realizadas em um repositório de vídeos
type VideoRepository interface {
	// Create persiste um novo vídeo no repositório
	// Retorna um erro se a operação falhar
	Create(ctx context.Context, video *entity.Video) error

	// FindByID busca um vídeo pelo seu ID
	// Retorna o vídeo encontrado ou um erro se não for encontrado ou se a operação falhar
	FindByID(ctx context.Context, id string) (*entity.Video, error)

	// List retorna uma lista de vídeos com paginação
	// page começa em 1, pageSize é o número de itens por página
	// Retorna a lista de vídeos ou um erro se a operação falhar
	List(ctx context.Context, page, pageSize int) ([]*entity.Video, error)

	// UpdateStatus atualiza o status de um vídeo e a mensagem de erro (quando aplicável)
	// Retorna um erro se a operação falhar
	UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error

	// UpdateHLSPath atualiza os caminhos HLS de um vídeo
	// Retorna um erro se a operação falhar
	UpdateHLSPath(ctx context.Context, id string, hlsPath, manifestPath string) error

	// UpdateS3Status atualiza o status de upload para S3 de um vídeo
	// Retorna um erro se a operação falhar
	UpdateS3Status(ctx context.Context, id string, uploadStatus string) error

	// UpdateS3URLs atualiza as URLs do S3 de um vídeo
	// Retorna um erro se a operação falhar
	UpdateS3URLs(ctx context.Context, id string, s3URL, s3ManifestURL string) error

	// UpdateS3Keys atualiza as chaves do S3 de um vídeo (segmentKey e manifestKey)
	// Retorna um erro se a operação falhar
	UpdateS3Keys(ctx context.Context, id string, segmentKey string, manifestKey string) error

	// Delete remove um vídeo do repositório
	// Retorna um erro se a operação falhar
	Delete(ctx context.Context, id string) error
}
