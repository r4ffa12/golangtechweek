package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devfullcycle/golangtechweek/internal/domain/entity"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// OutputFile representa um arquivo gerado pela conversão
type OutputFile struct {
	Path string // Caminho completo do arquivo
	Type string // Tipo do arquivo (manifest, segment)
}

// FFmpegServiceInterface define a interface para o serviço de conversão de vídeos usando FFmpeg.
// Esta interface facilita a criação de mocks para testes e segue o princípio de inversão de dependência.
//
// Exemplo de uso com mock em testes:
//
//	mockFFmpeg := new(MockFFmpegService)
//	mockFFmpeg.On("ConvertToHLS", ctx, inputPath, outputDir).Return(expectedFiles, nil)
//	// Use o mock no seu teste
type FFmpegServiceInterface interface {
	// ConvertToHLS converte um arquivo de vídeo para o formato HLS (HTTP Live Streaming).
	// Retorna uma lista de arquivos gerados (manifesto e segmentos) e um possível erro.
	ConvertToHLS(ctx context.Context, input string, outputDir string) ([]OutputFile, error)
}

// FFmpegService implementa a interface FFmpegServiceInterface usando o pacote ffmpeg-go.
// Esta estrutura não possui campos, pois não precisa armazenar estado.
type FFmpegService struct{}

// NewFFmpegService cria uma nova instância do serviço FFmpeg.
//
// Exemplo de uso:
//
//	ffmpegService := NewFFmpegService()
//	outputFiles, err := ffmpegService.ConvertToHLS(ctx, "video.mp4", "./output")
func NewFFmpegService() *FFmpegService {
	return &FFmpegService{}
}

// Sobre permissões de arquivos em notação octal (como 0o755):
//
// As permissões em sistemas Unix/Linux são representadas por 9 bits, organizados em 3 grupos:
// - Proprietário (dono do arquivo)
// - Grupo (usuários do mesmo grupo)
// - Outros (todos os demais usuários)
//
// Cada grupo tem 3 tipos de permissão:
// - r (read/leitura): Valor 4
// - w (write/escrita): Valor 2
// - x (execute/execução): Valor 1
//
// A notação octal é usada porque cada dígito octal representa exatamente 3 bits:
// - 7 (4+2+1) = rwx = leitura + escrita + execução
// - 5 (4+0+1) = r-x = leitura + execução
// - 6 (4+2+0) = rw- = leitura + escrita
// - 4 (4+0+0) = r-- = apenas leitura
//
// Exemplos comuns:
// - 0o755 (rwxr-xr-x): Diretórios padrão, scripts executáveis
// - 0o644 (rw-r--r--): Arquivos regulares (documentos, imagens)
// - 0o600 (rw-------): Arquivos privados (chaves SSH, senhas)
//
// No código abaixo, usamos 0o755 para o diretório de saída porque queremos que:
// 1. O processo que executa o código possa criar, modificar e listar arquivos
// 2. Outros usuários possam listar e acessar os arquivos, mas não modificar o diretório

// ConvertToHLS converte um vídeo para o formato HLS (HTTP Live Streaming).
// Parâmetros:
//   - ctx: Contexto que permite cancelamento da operação
//   - input: Caminho do arquivo de vídeo de entrada
//   - outputDir: Diretório onde os arquivos HLS serão salvos
//
// Retorna:
//   - Uma lista de arquivos gerados (manifesto e segmentos)
//   - Um erro, se ocorrer
//
// Exemplo de uso:
//
//	ctx := context.Background()
//	outputFiles, err := ffmpegService.ConvertToHLS(ctx, "video.mp4", "./output")
//	if err != nil {
//	    log.Fatalf("Erro ao converter vídeo: %v", err)
//	}
//	fmt.Printf("Arquivos gerados: %d\n", len(outputFiles))
func (s *FFmpegService) ConvertToHLS(ctx context.Context, input string, outputDir string) ([]OutputFile, error) {
	// Verifica se a operação já foi cancelada
	if ctx.Err() != nil {
		return nil, fmt.Errorf("operação cancelada: %w", ctx.Err())
	}

	// Cria o diretório de saída (e diretórios pai, se necessário)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de saída: %w", err)
	}

	// Executa a conversão do vídeo para o formato HLS
	if err := s.executeFFmpegConversion(ctx, input, outputDir); err != nil {
		return nil, fmt.Errorf("erro na conversão FFmpeg: %w", err)
	}

	// Coleta os arquivos gerados pela conversão
	return s.collectOutputFiles(outputDir)
}

// executeFFmpegConversion executa o comando FFmpeg para converter o vídeo para HLS.
// Esta função configura e executa o FFmpeg com os parâmetros necessários para
// criar um stream HLS a partir do vídeo de entrada.
func (s *FFmpegService) executeFFmpegConversion(ctx context.Context, input string, outputDir string) error {
	// Define o caminho do arquivo de manifesto (playlist principal)
	manifestPath := filepath.Join(outputDir, "playlist.m3u8")

	// Configura os parâmetros para a conversão HLS
	hlsParams := ffmpeg.KwArgs{
		// Parâmetros essenciais
		"f":             "hls", // Formato de saída: HLS
		"hls_time":      10,    // Duração de cada segmento em segundos
		"hls_list_size": 0,     // 0 = incluir todos os segmentos no manifesto

		// Codecs de áudio e vídeo
		"c:v": "h264", // Codec de vídeo H.264 (amplamente suportado)
		"c:a": "aac",  // Codec de áudio AAC (amplamente suportado)
		"b:a": "128k", // Taxa de bits do áudio: 128 kbps
	}

	// Executa o comando FFmpeg
	err := ffmpeg.Input(input).
		Output(manifestPath, hlsParams).
		ErrorToStdOut().
		Run()

	// Verifica se a operação foi cancelada durante a execução
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Retorna o erro do FFmpeg, se houver
	if err != nil {
		return err
	}

	return nil
}

// collectOutputFiles lista e categoriza os arquivos gerados pela conversão.
// Esta função percorre o diretório de saída e identifica os arquivos de manifesto (.m3u8)
// e os segmentos de vídeo (.ts) gerados pelo FFmpeg.
func (s *FFmpegService) collectOutputFiles(outputDir string) ([]OutputFile, error) {
	// Lista todos os arquivos no diretório de saída
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar arquivos gerados: %w", err)
	}

	// Cria um slice para armazenar os arquivos de saída
	outputFiles := make([]OutputFile, 0, len(files))

	// Processa cada arquivo encontrado
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(outputDir, fileName)

		// Determina o tipo do arquivo com base na extensão
		var fileType string
		if strings.HasSuffix(fileName, ".m3u8") {
			// Arquivos .m3u8 são manifestos (playlists)
			fileType = entity.FileTypeManifest
		} else {
			// Outros arquivos (geralmente .ts) são segmentos de vídeo
			fileType = entity.FileTypeSegment
		}

		// Adiciona o arquivo à lista de saída
		outputFiles = append(outputFiles, OutputFile{
			Path: filePath,
			Type: fileType,
		})
	}

	return outputFiles, nil
}
