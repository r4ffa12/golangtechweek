# Tarefa: Implementação da Entidade Video

## Descrição
Implementar a entidade de domínio Video, que será responsável por gerenciar o ciclo de vida completo do processamento de vídeos, desde o recebimento até a disponibilização no S3 no formato HLS.

## Objetivos
- [x] Criar a estrutura básica da entidade Video com todos os campos necessários
- [x] Definir as constantes para os possíveis estados do vídeo
- [x] Garantir que a entidade siga os princípios de Domain-Driven Design
- [x] Preparar a estrutura para futura implementação dos métodos
- [x] Implementar métodos para gerenciar o estado do vídeo
- [x] Implementar métodos para gerenciar URLs e caminhos de arquivos

## Campos Necessários
- [x] ID: Identificador único do vídeo
- [x] Title: Título do vídeo
- [x] FilePath: Caminho do arquivo original no sistema de arquivos
- [x] HLSPath: Caminho onde os arquivos HLS serão armazenados temporariamente
- [x] ManifestPath: Caminho do arquivo de manifesto (.m3u8)
- [x] S3URL: URL final do vídeo no S3 após o upload
- [x] S3ManifestURL: URL do manifesto no S3
- [x] Status: Estado atual do vídeo (pendente, em processamento, concluído, erro)
- [x] UploadStatus: Status do upload para o S3
- [x] ErrorMessage: Mensagem de erro, se houver
- [x] CreatedAt: Data de criação do registro
- [x] UpdatedAt: Data da última atualização do registro

## Métodos Implementados
- [x] NewVideo(title, filePath string): Construtor com geração automática de UUID
- [x] MarkAsProcessing(): Atualiza o status para "processing"
- [x] MarkAsCompleted(hslPath, manifestPath string): Atualiza o status para "completed"
- [x] MarkAsFailed(errorMessage string): Atualiza o status para "failed"
- [x] SetS3URL(url string): Define a URL final do vídeo no S3
- [x] SetS3ManifestURL(url string): Define a URL do manifesto no S3
- [x] IsCompleted(): Verifica se o vídeo foi processado com sucesso
- [x] GetHLSDirectory(): Retorna o diretório dos arquivos HLS
- [x] GetManifestPath(): Retorna o caminho do manifesto
- [x] GenerateOutputPath(baseDir string): Gera o caminho base para os arquivos convertidos

## Estrutura de Diretórios
- [x] Criar diretório para o domínio (internal/domain)
- [x] Criar arquivo para a entidade Video (internal/domain/entity/video.go)
- [x] Criar arquivo para testes da entidade (internal/domain/entity/video_test.go)