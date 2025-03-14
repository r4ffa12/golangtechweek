package entity

import (
	"time"

	"github.com/google/uuid"
)

// Status do vídeo durante o ciclo de processamento
const (
	// StatusPending representa um vídeo que foi registrado mas ainda não começou a ser processado
	StatusPending = "pending"
	// StatusProcessing representa um vídeo que está sendo processado
	StatusProcessing = "processing"

	// StatusCompleted representa um vídeo que foi processado com sucesso
	StatusCompleted = "completed"

	// StatusError representa um vídeo que encontrou um erro durante o processamento
	StatusError = "failed"
)

const (
	UploadStatusNone        = "none"
	UploadStatusPendingS3   = "pending_s3"
	UploadStatusUploadingS3 = "uploading_s3"
	UploadStatusCompletedS3 = "completed_s3"
	UploadStatusFailedS3    = "failed_s3"
)

const (
	FileTypeManifest = "manifest"
	FileTypeSegment  = "segment"
)

// Video representa a entidade de domínio para um vídeo que será processado
type Video struct {
	ID            string // Identificador único do vídeo
	Title         string // Título do vídeo
	FilePath      string // Caminho do arquivo original no sistema de arquivos
	HLSPath       string // Caminho onde os arquivos HLS serão armazenados temporariamente
	ManifestPath  string // Caminho do arquivo de manifesto (.m3u8)
	S3ManifestURL string // URL do manifesto no S3
	S3URL         string // URL final do vídeo no S3 após o upload
	Status        string // Estado atual do vídeo
	UploadStatus  string
	ErrorMessage  string    // Mensagem de erro, se houver
	CreatedAt     time.Time // Data de criação do registro
	UpdatedAt     time.Time // Data da última atualização do registro
}

// NewVideo cria uma nova instância de Video com valores padrão
func NewVideo(title, filePath string) *Video {
	now := time.Now()

	return &Video{
		ID:           uuid.New().String(),
		Title:        title,
		FilePath:     filePath,
		Status:       StatusPending,
		UploadStatus: UploadStatusNone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// MarkAsProcessing atualiza o status do vídeo para "processing"
func (v *Video) MarkAsProcessing() {
	v.Status = StatusProcessing
	v.UpdatedAt = time.Now()
}

// MarkAsCompleted atualiza o status do vídeo para "completed"
func (v *Video) MarkAsCompleted(hslPath, manifestPath string) {
	v.Status = StatusCompleted
	v.HLSPath = hslPath
	v.ManifestPath = manifestPath
	v.UpdatedAt = time.Now()
}

// MarkAsFailed atualiza o status do vídeo para "failed" e registra a mensagem de erro
func (v *Video) MarkAsFailed(errorMessage string) {
	v.Status = StatusError
	v.ErrorMessage = errorMessage
	v.UpdatedAt = time.Now()
}

// SetS3URL define a URL final do vídeo no S3
func (v *Video) SetS3URL(url string) {
	v.S3URL = url
	v.UpdatedAt = time.Now()
}

// SetS3ManifestURL define a URL do manifesto no S3
func (v *Video) SetS3ManifestURL(url string) {
	v.S3ManifestURL = url
	v.UpdatedAt = time.Now()
}

// IsCompleted verifica se o vídeo foi processado com sucesso
func (v *Video) IsCompleted() bool {
	return v.Status == StatusCompleted
}

// GetHLSDirectory retorna o diretório onde os arquivos HLS estão armazenados
func (v *Video) GetHLSDirectory() string {
	return v.HLSPath
}

// GetManifestPath retorna o caminho do arquivo de manifesto
func (v *Video) GetManifestPath() string {
	return v.ManifestPath
}

// GenerateOutputPath gera os caminhos de saída para os arquivos HLS
func (v *Video) GenerateOutputPath(baseDir string) string {
	return baseDir + "/converted/" + v.ID
}
