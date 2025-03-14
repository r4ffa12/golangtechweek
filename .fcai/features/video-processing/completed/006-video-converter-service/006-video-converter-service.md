# Tarefa 006: Implementação do Serviço de Conversão de Vídeo com Worker Pool

## Descrição
Implementar um serviço de conversão de vídeo que utilize o worker pool para processar múltiplas solicitações de conversão simultaneamente. O serviço deve atualizar o status do vídeo no banco de dados conforme o processamento avança.

## Objetivos
- [x] Definir as estruturas `ConversionJob` e `ConversionResult`
- [x] Implementar o serviço `VideoConverterService` com as dependências necessárias
- [x] Implementar a função de processamento `processFunc` para o worker pool
- [x] Implementar o método `StartConversion` para iniciar o processo de conversão
- [x] Implementar a atualização de status do vídeo no banco de dados
- [ ] Criar testes unitários e de integração para o serviço

## Especificações Técnicas

### Estrutura do Job de Conversão
```go
type ConversionJob struct {
    VideoID   string
    InputPath string
    OutputDir string
}
```

### Estrutura do Resultado da Conversão
```go
type ConversionResult struct {
    VideoID     string
    Success     bool
    Error       error
    OutputFiles []service.OutputFile
    Duration    time.Duration
}
```

### Interface do Serviço
```go
type VideoConverterService interface {
    StartConversion(ctx context.Context, inputCh <-chan ConversionJob) (<-chan ConversionResult, error)
    Stop() error
    IsRunning() bool
}
```

### Dependências do Serviço
- FFmpegService: Para realizar a conversão de vídeos
- VideoRepository: Para atualizar o status dos vídeos no banco de dados
- WorkerPool: Para processar múltiplas conversões simultaneamente
- Logger: Para registrar eventos e erros durante o processamento

### Localização dos Arquivos
- **Serviço**: `internal/application/service/video_converter.go`
- **Testes**: `internal/application/service/video_converter_test.go`

## Critérios de Aceitação
1. O serviço deve processar múltiplas solicitações de conversão simultaneamente
2. O status do vídeo deve ser atualizado no banco de dados conforme o processamento
3. O serviço deve retornar os resultados da conversão através de um canal
4. Os testes devem validar o funcionamento correto do serviço com mocks para as dependências

## Dependências
- Serviço FFmpeg implementado (Tarefa 005)
- Repositório de vídeos implementado (Tarefa 002/003)
- Worker Pool implementado (pkg/workerpool)

## Notas Adicionais
- O serviço deve ser implementado seguindo os princípios de Clean Architecture
- A implementação deve garantir o tratamento adequado de erros durante a conversão
- O serviço deve ser configurável quanto ao número de workers utilizados 