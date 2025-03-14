# Tarefa: Implementação da Interface de Repositório de Vídeos

## Descrição
Implementar a interface de repositório de vídeos dentro do domínio, seguindo os princípios de Clean Architecture e Domain-Driven Design. Esta interface será responsável por definir as operações que podem ser realizadas em um repositório de vídeos, permitindo diferentes implementações (PostgreSQL, MongoDB, em memória, etc.) sem alterar o domínio.

## Objetivos
- [x] Criar a estrutura de diretórios para o repositório
- [x] Definir a interface VideoRepository com todos os métodos necessários
- [x] Garantir que a interface siga os princípios de Clean Architecture
- [x] Preparar a estrutura para futura implementação concreta

## Métodos Necessários
- [x] Create(ctx context.Context, video *entity.Video) error
- [x] FindByID(ctx context.Context, id string) (*entity.Video, error)
- [x] List(ctx context.Context, page, pageSize int) ([]*entity.Video, error)
- [x] UpdateStatus(ctx context.Context, id string, status string, errorMessage string) error
- [x] UpdateHLSPath(ctx context.Context, id string, hlsPath, manifestPath string) error
- [x] UpdateS3Status(ctx context.Context, id string, uploadStatus string) error
- [x] UpdateS3URLs(ctx context.Context, id string, s3URL, s3ManifestURL string) error
- [x] UpdateS3Keys(ctx context.Context, id string, segmentKey string, manifestKey string) error
- [x] Delete(ctx context.Context, id string) error

## Considerações de Design
- Todos os métodos devem receber um `context.Context` como primeiro parâmetro
- O método `Create` não deve retornar a entidade, apenas um erro se a operação falhar
- O método `List` deve incluir parâmetros de paginação (page e pageSize) e retornar apenas a lista de vídeos e um erro
- O método `UpdateStatus` deve incluir um parâmetro para a mensagem de erro, útil quando o status é "failed"
- Métodos de atualização retornam apenas `error`, não a entidade atualizada
- O método `UpdateS3Keys` recebe `segmentKey` e `manifestKey` em vez de uma lista de chaves
- Métodos específicos para cada tipo de atualização em vez de um único método `Update`

## Estrutura de Diretórios
- [x] Criar diretório para o repositório (internal/domain/repository)
- [x] Criar arquivo para a interface VideoRepository (internal/domain/repository/video_repository.go)

## Critérios de Aceitação
- A interface deve seguir os princípios de Clean Architecture
- A interface deve ser clara e fácil de implementar
- A interface deve permitir diferentes implementações sem alterar o domínio
- A documentação dos métodos deve ser clara e completa

## Dependências
- Entidade Video (já implementada)

## Estimativa
- 1 hora

## Responsável
- Equipe de desenvolvimento 