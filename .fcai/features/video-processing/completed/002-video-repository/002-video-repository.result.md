# Resultado: Implementação da Interface de Repositório de Vídeos

## Resumo
A implementação da interface de repositório de vídeos foi concluída com sucesso. Criamos a interface `VideoRepository` que define todas as operações necessárias para persistir e recuperar vídeos, seguindo os princípios de Clean Architecture e Domain-Driven Design.

## Implementação

### Estrutura de Diretórios
Criamos a seguinte estrutura de diretórios:
```
internal/
└── domain/
    └── repository/
        └── video_repository.go
```

### Interface VideoRepository
A interface `VideoRepository` foi implementada com os seguintes métodos:

```go
type VideoRepository interface {
	Create(ctx context.Context, video *entity.Video) error
	FindByID(ctx context.Context, id string) (*entity.Video, error)
	List(ctx context.Context, page, pageSize int) ([]*entity.Video, error)
	UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error
	UpdateHLSPath(ctx context.Context, id string, hlsPath, manifestPath string) error
	UpdateS3Status(ctx context.Context, id string, uploadStatus string) error
	UpdateS3URLs(ctx context.Context, id string, s3URL, s3ManifestURL string) error
	UpdateS3Keys(ctx context.Context, id string, segmentKey string, manifestKey string) error
	Delete(ctx context.Context, id string) error
}
```

### Características da Interface

1. **Uso de Context**: Todos os métodos recebem um `context.Context` como primeiro parâmetro, permitindo controle de timeout, cancelamento e passagem de valores entre camadas.

2. **Métodos de Criação e Consulta**:
   - `Create`: Persiste um novo vídeo e retorna apenas um erro se a operação falhar
   - `FindByID`: Busca um vídeo pelo ID e retorna a entidade encontrada ou um erro
   - `List`: Retorna uma lista paginada de vídeos

3. **Métodos de Atualização**:
   - Todos os métodos de atualização retornam apenas `error`, não a entidade atualizada
   - `UpdateStatus`: Inclui um parâmetro para a mensagem de erro, útil quando o status é "failed"
   - `UpdateHLSPath`: Atualiza os caminhos dos arquivos HLS
   - `UpdateS3Status`: Atualiza o status de upload para S3
   - `UpdateS3URLs`: Atualiza as URLs do S3
   - `UpdateS3Keys`: Atualiza as chaves do S3 (segmentKey e manifestKey)

4. **Método de Remoção**:
   - `Delete`: Remove um vídeo do repositório

### Considerações de Design

- **Separação de Responsabilidades**: Cada método tem uma responsabilidade clara e específica
- **Métodos Específicos**: Em vez de um único método `Update`, temos métodos específicos para cada tipo de atualização
- **Simplicidade**: Os métodos de atualização retornam apenas `error`, simplificando a interface
- **Flexibilidade**: A interface permite diferentes implementações (PostgreSQL, MongoDB, em memória, etc.) sem alterar o domínio

## Próximos Passos
- Implementar a interface em uma camada de infraestrutura (por exemplo, PostgreSQL)
- Criar testes para a implementação
- Injetar a implementação nos serviços que precisam dela 