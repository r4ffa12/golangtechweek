# Feature: Processamento de Vídeo

## Descrição
Esta feature é responsável por todo o ciclo de vida do processamento de vídeos, desde o recebimento do arquivo até a disponibilização do vídeo convertido para HLS no S3.

## Componentes Principais

### 1. Entidade de Domínio: Video
Representa o vídeo e seu estado durante todo o processo de conversão e upload. Contém todas as informações necessárias para rastrear e gerenciar o ciclo de vida do vídeo.

### 2. Serviço de Conversão
Responsável por converter o vídeo para o formato HLS utilizando o ffmpeg.

### 3. Serviço de Upload
Responsável por fazer o upload dos segmentos HLS e do manifesto para o AWS S3.

### 4. API REST
Endpoints para upload de vídeos, consulta de status e obtenção da URL final.

## Fluxo de Processamento
1. Recebimento do vídeo via API
2. Armazenamento temporário no servidor
3. Registro no banco de dados
4. Conversão para HLS (processamento concorrente)
5. Upload para S3
6. Atualização do status e URL final
7. Disponibilização do link para o cliente

## Tecnologias Utilizadas
- Go 1.24
- ffmpeg para conversão
- AWS SDK para Go
- PostgreSQL para persistência
- Goroutines e channels para concorrência 