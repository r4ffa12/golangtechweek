# Resultado: Implementação da Entidade Video

## Resumo
A implementação da entidade Video foi concluída com sucesso. Criamos a estrutura básica da entidade com todos os campos necessários, definimos as constantes para os possíveis estados do vídeo e implementamos os métodos para gerenciar o estado do vídeo durante o processamento.

## Implementação

### Estrutura de Diretórios
Criamos a seguinte estrutura de diretórios:
```
internal/
└── domain/
    └── entity/
        ├── video.go
        └── video_test.go
```

### Entidade Video
A entidade Video foi implementada com os seguintes campos:
- ID: Identificador único do vídeo (gerado automaticamente como UUID)
- Title: Título do vídeo
- FilePath: Caminho do arquivo original no sistema de arquivos
- HLSPath: Caminho onde os arquivos HLS serão armazenados temporariamente
- ManifestPath: Caminho do arquivo de manifesto (.m3u8)
- S3ManifestURL: URL do manifesto no S3
- S3URL: URL final do vídeo no S3 após o upload
- Status: Estado atual do vídeo
- UploadStatus: Status do upload para o S3
- ErrorMessage: Mensagem de erro, se houver
- CreatedAt: Data de criação do registro
- UpdatedAt: Data da última atualização do registro

### Constantes de Status
Definimos as seguintes constantes para os possíveis estados do vídeo:
- StatusPending: Vídeo registrado mas ainda não processado
- StatusProcessing: Vídeo em processamento
- StatusCompleted: Vídeo processado com sucesso
- StatusError: Vídeo com erro durante o processamento (definido como "failed")

### Constantes de Status de Upload
Definimos as seguintes constantes para os possíveis estados de upload:
- UploadStatusNone: Nenhum upload iniciado
- UploadStatusPendingS3: Upload pendente para o S3
- UploadStatusUploadingS3: Upload em andamento para o S3
- UploadStatusCompletedS3: Upload concluído para o S3
- UploadStatusFailedS3: Erro durante o upload para o S3

### Constantes de Tipo de Arquivo
Definimos as seguintes constantes para os tipos de arquivo:
- FileTypeManifest: Arquivo de manifesto (.m3u8)
- FileTypeSegment: Segmento de vídeo (.ts)

### Método Construtor
Implementamos o método `NewVideo` que cria uma nova instância de Video com valores padrão:
- Recebe apenas o título e o caminho do arquivo como parâmetros
- Gera automaticamente um UUID como ID
- Define o status inicial como "pending" e o status de upload como "none"
- Inicializa as datas de criação e atualização

### Métodos de Gerenciamento de Estado
Implementamos os seguintes métodos para gerenciar o estado do vídeo:
- `MarkAsProcessing()`: Atualiza o status do vídeo para "processing"
- `MarkAsCompleted(hslPath, manifestPath string)`: Atualiza o status do vídeo para "completed" e define os caminhos dos arquivos HLS
- `MarkAsFailed(errorMessage string)`: Atualiza o status do vídeo para "failed" e registra a mensagem de erro
- `IsCompleted()`: Verifica se o vídeo foi processado com sucesso

### Métodos de Gerenciamento de URLs
Implementamos os seguintes métodos para gerenciar as URLs do vídeo:
- `SetS3URL(url string)`: Define a URL final do vídeo no S3
- `SetS3ManifestURL(url string)`: Define a URL do manifesto no S3

### Métodos de Gerenciamento de Caminhos
Implementamos os seguintes métodos para gerenciar os caminhos dos arquivos:
- `GetHLSDirectory()`: Retorna o diretório onde os arquivos HLS estão armazenados
- `GetManifestPath()`: Retorna o caminho do arquivo de manifesto
- `GenerateOutputPath(baseDir string) string`: Gera e retorna o caminho base para os arquivos convertidos, seguindo o padrão `baseDir/converted/videoID`

### Modificações Importantes
1. O método `MarkAsCompleted` agora recebe os parâmetros `hslPath` e `manifestPath` para definir os caminhos dos arquivos HLS no momento em que o vídeo é marcado como concluído.

2. O método `GenerateOutputPath` foi modificado para:
   - Retornar uma string com o caminho de saída
   - Incluir um subdiretório "converted" na estrutura do caminho
   - Não atualizar diretamente os campos da entidade, apenas gerar e retornar o caminho

Essas modificações tornam a API da entidade mais flexível, permitindo que o serviço de conversão tenha mais controle sobre os caminhos dos arquivos e quando eles são definidos na entidade.

### Testes
Implementamos testes unitários para verificar:
- Se o ID é gerado automaticamente
- Se o construtor inicializa corretamente os campos obrigatórios
- Se os status iniciais estão corretos
- Se as datas são preenchidas corretamente
- Se os campos opcionais estão vazios inicialmente
- Se os métodos de gerenciamento de estado funcionam corretamente
- Se os métodos de gerenciamento de URLs funcionam corretamente
- Se os métodos de gerenciamento de caminhos funcionam corretamente

## Próximos Passos
- Implementar o serviço de conversão para HLS
- Implementar o serviço de upload para S3
- Implementar a API REST para upload de vídeos 