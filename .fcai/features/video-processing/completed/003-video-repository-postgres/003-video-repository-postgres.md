# Tarefa: Implementação do Repositório PostgreSQL para Vídeos

## Descrição
Implementar a versão concreta do repositório de vídeos utilizando PostgreSQL, seguindo a interface definida em `internal/domain/repository/video_repository.go`. Esta implementação será responsável por persistir e recuperar vídeos do banco de dados PostgreSQL.

## Objetivos
- [x] Criar a estrutura de migrations para o banco de dados
- [x] Implementar a migration para criação da tabela de vídeos
- [x] Implementar a conexão com o banco de dados PostgreSQL
- [x] Implementar a versão concreta do repositório de vídeos
- [x] Implementar testes de integração para o repositório

## Subtarefas

### 1. Estrutura de Migrations
- [x] Criar diretório para migrations (`internal/infra/database/migrations`)
- [x] Configurar o pacote `golang-migrate/migrate` para gerenciar as migrations
- [x] Usar o `migrate` para criar as migrations

### 2. Criação da Tabela de Vídeos
- [x] Criar migration para criação da tabela `videos` com os seguintes campos:
```
CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_path VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    upload_status VARCHAR(50) NOT NULL DEFAULT 'none',
    error_message TEXT,
    hls_path VARCHAR(255),
    manifest_path VARCHAR(255),
    s3_url VARCHAR(255),
    s3_manifest_url VARCHAR(255),
    segment_key VARCHAR(255),
    manifest_key VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
)
```
- [x] Criar migration para rollback (drop table)

### 3. Implementação do Repositório
- [x] Criar estrutura de diretórios (`internal/infra/database/repository`)
- [x] Implementar a conexão com o banco de dados
- [x] Implementar o método `Create`
- [x] Implementar o método `FindByID`
- [x] Implementar o método `List`
- [x] Implementar o método `UpdateStatus`
- [x] Implementar o método `UpdateHLSPath`
- [x] Implementar o método `UpdateS3Status`
- [x] Implementar o método `UpdateS3URLs`
- [x] Implementar o método `UpdateS3Keys`
- [x] Implementar o método `Delete`

### 4. Testes de Integração
- [x] Configurar ambiente de teste com banco de dados de teste
- [x] Implementar testes para o método `Create`
- [x] Implementar testes para o método `FindByID`
- [x] Implementar testes para o método `List`
- [x] Implementar testes para os métodos de atualização
- [x] Implementar testes para o método `Delete`

## Estrutura de Diretórios
- [x] `internal/infra/database/migrations` - Migrations do banco de dados
- [x] `internal/infra/database/repository` - Implementação do repositório PostgreSQL
- [x] `internal/infra/database/repository/video_repository_test.go` - Testes de integração

## Critérios de Aceitação
- A implementação deve seguir a interface definida em `internal/domain/repository/video_repository.go`
- As migrations devem criar corretamente a tabela de vídeos com todos os campos necessários
- Os testes de integração devem verificar o funcionamento correto de todos os métodos
- O código deve seguir as boas práticas de programação em Go
- A implementação deve ser thread-safe e lidar corretamente com erros de banco de dados

## Dependências
- Interface VideoRepository (já implementada)
- Entidade Video (já implementada)
- PostgreSQL (disponível via Docker)

## Pacotes Necessários
- `database/sql` - Pacote padrão do Go para SQL
- `github.com/lib/pq` - Driver PostgreSQL para Go
- `github.com/golang-migrate/migrate/v4` - Ferramenta para migrations
