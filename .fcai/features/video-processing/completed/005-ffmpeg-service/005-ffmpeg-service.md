# Tarefa 005: Implementação do Serviço de Conversão de Vídeo (FFmpeg)

## Descrição
Implementar um serviço para conversão de vídeos para o formato HLS (HTTP Live Streaming) utilizando FFmpeg. O serviço deve fornecer uma interface clara para converter vídeos e coletar os arquivos gerados.

## Objetivos
- [x] Criar a interface `FFmpegService` com o método `ConvertToHLS`
- [x] Implementar a estrutura `FFmpegConverter` que implementa a interface
- [x] Implementar o método de execução do FFmpeg com parâmetros configuráveis
- [x] Implementar o método para coletar os arquivos gerados
- [x] Criar testes de integração para validar a conversão

## Especificações Técnicas

### Interface do Serviço
```go
type FFmpegService interface {
    ConvertToHLS(ctx context.Context, input string, outputDir string) ([]OutputFile, error)
}
```

### Estrutura do Conversor
```go
type FFmpegConverter struct {
    videoCodec   string
    audioCodec   string
    videoBitrate string
    audioBitrate string
}
```

### Estrutura de Saída
```go
type OutputFile struct {
    Path string // Caminho completo do arquivo
    Type string // Tipo do arquivo (manifest, segment)
}
```

### Parâmetros de Configuração
- **videoCodec**: Codec de vídeo a ser utilizado (ex: "libx264")
- **audioCodec**: Codec de áudio a ser utilizado (ex: "aac")
- **videoBitrate**: Taxa de bits para o vídeo (ex: "1000k")
- **audioBitrate**: Taxa de bits para o áudio (ex: "128k")

### Localização dos Arquivos
- **Serviço**: `internal/application/service/ffmpeg_service.go`
- **Testes**: `internal/application/service/ffmpeg_service_test.go`

## Critérios de Aceitação
1. O serviço deve converter com sucesso um arquivo de vídeo para o formato HLS
2. Os arquivos gerados devem incluir um manifesto (.m3u8) e segmentos (.ts)
3. O serviço deve retornar a lista de arquivos gerados com seus respectivos tipos
4. Os testes de integração devem passar utilizando um arquivo de vídeo real

## Dependências
- FFmpeg instalado no ambiente de execução
- Estrutura de diretórios para armazenar os arquivos temporários

## Notas Adicionais
- Para os testes de integração, serão utilizados arquivos de vídeo de teste localizados na pasta `upload/`
- O serviço deve ser implementado de forma a facilitar a integração com o serviço de upload para S3 que será implementado posteriormente 