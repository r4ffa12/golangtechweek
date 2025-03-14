# Estrutura de Camadas e Responsabilidades do ConversorGo

## Visão Geral

O ConversorGo segue uma arquitetura em camadas para organizar seu código, garantindo separação de responsabilidades, testabilidade e manutenibilidade. Este documento descreve a estrutura arquitetural do projeto, as responsabilidades de cada camada e as diretrizes para implementação.

## Princípios Fundamentais

1. **Separação de Responsabilidades**: Cada camada tem um propósito específico e bem definido
2. **Testabilidade**: Todas as camadas podem ser testadas de forma isolada
3. **Independência de Detalhes Técnicos**: O núcleo da aplicação não depende de detalhes de implementação
4. **Inversão de Dependência**: Dependências apontam para dentro, não para fora
5. **Substituibilidade**: Componentes podem ser substituídos sem afetar o restante do sistema

## Estrutura de Camadas

O ConversorGo é organizado nas seguintes camadas:

```
internal/
├── domain/             # Regras de negócio e entidades
├── application/        # Casos de uso e orquestração de serviços
└── infra/     # Implementações técnicas e interfaces com o mundo externo
```

### 1. Domain (Domínio)

A camada de domínio contém as entidades e regras de negócio centrais da aplicação, independentes de qualquer detalhe de implementação.

#### Responsabilidades:
- Definir entidades de negócio
- Definir interfaces de repositórios
- Implementar regras de negócio puras

#### Estrutura:
```
domain/
├── entity/         # Entidades de negócio
└── repository/     # Interfaces de repositórios
```

#### Diretrizes:
- Não deve depender de nenhuma outra camada
- Não deve importar pacotes externos exceto os da biblioteca padrão Go
- Deve conter apenas regras de negócio puras
- Deve definir interfaces que serão implementadas por camadas externas

#### Exemplo:
```go
// domain/entity/video.go
package entity

type Video struct {
    ID           string
    Title        string
    FilePath     string
    Status       string
    // ...
}

func (v *Video) CanBeProcessed() bool {
    return v.Status == "pending" || v.Status == "failed"
}

// domain/repository/video_repository.go
package repository

type VideoRepository interface {
    FindByID(ctx context.Context, id string) (*entity.Video, error)
    UpdateStatus(ctx context.Context, id, status, errorMessage string) error
    // ...
}
```

### 2. Application (Aplicação)

A camada de aplicação contém os casos de uso da aplicação, orquestrando o fluxo entre entidades e serviços de infraestrutura.

#### Responsabilidades:
- Implementar casos de uso
- Orquestrar entre domínio e infraestrutura
- Gerenciar transações e fluxo de dados
- Implementar serviços de aplicação

#### Estrutura:
```
application/
├── service/        # Serviços de aplicação
└── usecase/        # Casos de uso específicos
```

#### Diretrizes:
- Pode depender apenas da camada de domínio
- Não deve conter regras de negócio, apenas orquestração
- Pode receber implementações de infraestrutura via injeção de dependência

#### Exemplo:
```go
// application/service/video_converter.go
package service

import (
    "github.com/devfullcycle/golangtechweek/internal/domain/repository"
)

// VideoConverter é um serviço de aplicação para conversão de vídeos
type VideoConverter struct {
    ffmpeg      FFmpegWrapper
    videoRepo   repository.VideoRepository
    // ...
}

// Interface para o wrapper FFmpeg (definida na aplicação)
type FFmpegWrapper interface {
    ConvertToHLS(ctx context.Context, input, outputDir string) ([]OutputFile, error)
    GetVideoInfo(ctx context.Context, input string) (VideoInfo, error)
}

// StartConversion inicia o serviço de conversão
func (c *VideoConverter) StartConversion(ctx context.Context, inputCh <-chan ConversionJob) (<-chan ConversionResult, error) {
    // Implementação...
}
```

### 3. Infrastructure (infra)

A camada de infraestrutura contém implementações concretas de interfaces definidas no domínio e na aplicação, bem como todos os componentes que interagem com o mundo externo.

#### Responsabilidades:
- Implementar repositórios
- Integrar com serviços externos
- Fornecer adaptadores para frameworks
- Implementar detalhes técnicos
- Fornecer interfaces com o mundo externo (API, CLI)

#### Estrutura:
```
infra/
├── database/       # Implementações de banco de dados
├── repository/     # Implementações de repositórios
├── s3/        # Serviços de armazenamento
├── api/            # API HTTP
│   ├── handler/    # Handlers HTTP
│   └── router.go   # Configuração de rotas
```

#### Diretrizes:
- Pode depender das camadas de domínio e aplicação
- Deve implementar interfaces definidas no domínio e na aplicação
- Deve encapsular detalhes técnicos
- Deve ser substituível sem afetar as camadas internas



## Fluxo de Dependências

O fluxo de dependências segue a regra da dependência: as camadas internas não conhecem as camadas externas.

```
Infrastructure → Application → Domain
```

## Injeção de Dependências

A injeção de dependências é utilizada para fornecer implementações concretas para interfaces:

```go
// cmd/app/main.go
func main() {
    // Infraestrutura
    db := postgres.NewConnection()
    videoRepo := postgres.NewVideoRepository(db)
    ffmpegWrapper := ffmpeg.NewFFmpegWrapper()
    
    // Aplicação
    videoConverter := service.NewVideoConverter(ffmpegWrapper, videoRepo)
    uploadUseCase := usecase.NewUploadVideo(videoRepo, videoConverter)
    
    // API (parte da infraestrutura)
    videoHandler := handler.NewVideoHandler(uploadUseCase)
    router := api.NewRouter(videoHandler)
    
    // Iniciar servidor
    http.ListenAndServe(":8080", router)
}
```
