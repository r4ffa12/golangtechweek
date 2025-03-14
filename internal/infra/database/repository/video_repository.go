package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	domainRepository "github.com/devfullcycle/golangtechweek/internal/domain/repository"
)

var ErrVideoNotFound = errors.New("vídeo não encontrado")

// VideoRepositoryPostgres implementa a interface VideoRepository usando PostgreSQL
type VideoRepositoryPostgres struct {
	db *sql.DB
}

// NewVideoRepositoryPostgres cria uma nova instância de VideoRepositoryPostgres
func NewVideoRepositoryPostgres(db *sql.DB) *VideoRepositoryPostgres {
	return &VideoRepositoryPostgres{
		db: db,
	}
}

// Create persiste um novo vídeo no banco de dados
func (r *VideoRepositoryPostgres) Create(ctx context.Context, video *entity.Video) error {
	query := `
		INSERT INTO videos (
			id, title, file_path, status, upload_status, hls_path, manifest_path, 
			s3_url, s3_manifest_url, error_message, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		video.ID,
		video.Title,
		video.FilePath,
		video.Status,
		video.UploadStatus,
		video.HLSPath,
		video.ManifestPath,
		video.S3URL,
		video.S3ManifestURL,
		video.ErrorMessage,
		video.CreatedAt,
		video.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar vídeo: %w", err)
	}

	return nil
}

// FindByID busca um vídeo pelo seu ID
func (r *VideoRepositoryPostgres) FindByID(ctx context.Context, id string) (*entity.Video, error) {
	query := `
		SELECT 
			id, title, file_path, status, upload_status, hls_path, manifest_path, 
			s3_url, s3_manifest_url, error_message, created_at, updated_at
		FROM videos
		WHERE id = $1 AND deleted_at IS NULL
	`

	var video entity.Video
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&video.ID,
		&video.Title,
		&video.FilePath,
		&video.Status,
		&video.UploadStatus,
		&video.HLSPath,
		&video.ManifestPath,
		&video.S3URL,
		&video.S3ManifestURL,
		&video.ErrorMessage,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVideoNotFound
		}
		return nil, fmt.Errorf("erro ao buscar vídeo: %w", err)
	}

	video.CreatedAt = createdAt
	video.UpdatedAt = updatedAt

	return &video, nil
}

// List retorna uma lista de vídeos com paginação
func (r *VideoRepositoryPostgres) List(ctx context.Context, page, pageSize int) ([]*entity.Video, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	query := `
		SELECT 
			id, title, file_path, status, upload_status, hls_path, manifest_path, 
			s3_url, s3_manifest_url, error_message, created_at, updated_at
		FROM videos
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar vídeos: %w", err)
	}
	defer rows.Close()

	var videos []*entity.Video

	for rows.Next() {
		var video entity.Video
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&video.ID,
			&video.Title,
			&video.FilePath,
			&video.Status,
			&video.UploadStatus,
			&video.HLSPath,
			&video.ManifestPath,
			&video.S3URL,
			&video.S3ManifestURL,
			&video.ErrorMessage,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear vídeo: %w", err)
		}

		video.CreatedAt = createdAt
		video.UpdatedAt = updatedAt

		videos = append(videos, &video)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", err)
	}

	return videos, nil
}

// UpdateStatus atualiza o status de um vídeo e a mensagem de erro (quando aplicável)
func (r *VideoRepositoryPostgres) UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error {
	query := `
		UPDATE videos
		SET status = $1, error_message = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, status, errorMessage, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// UpdateHLSPath atualiza os caminhos HLS de um vídeo
func (r *VideoRepositoryPostgres) UpdateHLSPath(ctx context.Context, id string, hlsPath, manifestPath string) error {
	query := `
		UPDATE videos
		SET hls_path = $1, manifest_path = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, hlsPath, manifestPath, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar caminhos HLS do vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// UpdateS3Status atualiza o status de upload para S3 de um vídeo
func (r *VideoRepositoryPostgres) UpdateS3Status(ctx context.Context, id string, uploadStatus string) error {
	query := `
		UPDATE videos
		SET upload_status = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, uploadStatus, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status de upload do vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// UpdateS3URLs atualiza as URLs do S3 de um vídeo
func (r *VideoRepositoryPostgres) UpdateS3URLs(ctx context.Context, id string, s3URL, s3ManifestURL string) error {
	query := `
		UPDATE videos
		SET s3_url = $1, s3_manifest_url = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, s3URL, s3ManifestURL, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar URLs do S3 do vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// UpdateS3Keys atualiza as chaves do S3 de um vídeo
func (r *VideoRepositoryPostgres) UpdateS3Keys(ctx context.Context, id string, segmentKey string, manifestKey string) error {
	query := `
		UPDATE videos
		SET segment_key = $1, manifest_key = $2, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, segmentKey, manifestKey, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar chaves do S3 do vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// Delete remove um vídeo do repositório (soft delete)
func (r *VideoRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE videos
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("erro ao excluir vídeo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao obter linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return ErrVideoNotFound
	}

	return nil
}

// Ensure VideoRepositoryPostgres implements VideoRepository
var _ domainRepository.VideoRepository = (*VideoRepositoryPostgres)(nil)
