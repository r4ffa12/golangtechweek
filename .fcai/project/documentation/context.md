# Contexto Geral do Projeto: Conversão de Vídeos para HLS com Upload para S3

## **Descrição do Projeto**
O projeto tem como objetivo criar uma **API REST escalável e eficiente** para processamento de vídeos, convertendo-os para o formato **HLS (HTTP Live Streaming)** e armazenando os segmentos no **AWS S3**. A aplicação será desenvolvida em **Go 1.24**, utilizando um modelo baseado em **concorrência eficiente com goroutines e channels** para processar e gerenciar múltiplas conversões simultaneamente.

A solução permitirá que aplicações externas enviem vídeos via **API HTTP**, iniciando automaticamente a conversão e disponibilizando os arquivos no S3. Os usuários poderão consultar o status da conversão e obter a URL final para visualização do vídeo convertido.

## **Objetivo do Projeto**
- Fornecer um serviço de conversão de vídeos **altamente eficiente** e **escalável**.
- Utilizar **concorrência em Go** para permitir múltiplas conversões simultâneas sem impactar o desempenho do sistema.
- Implementar uma API **simples e intuitiva** para upload de vídeos e consulta de status.
- Armazenar os vídeos processados em **AWS S3**, otimizando o acesso e distribuição dos arquivos convertidos. Utilizar o Localstack para testar o upload para S3 localmente.

## **Principais Tecnologias Utilizadas**
- **Linguagem:** Go 1.24
- **Processamento de Mídia:** `ffmpeg` para conversão dos vídeos para HLS
- **Banco de Dados:** PostgreSQL para rastrear status e metadados das conversões
- **Armazenamento:** AWS S3 para guardar os segmentos `.ts` e os manifestos `.m3u8`
- **Concorrência:** Goroutines e Channels para gerenciar múltiplas conversões simultaneamente
- **APIs:** RESTful HTTP para upload, consulta de status e obtenção do link final do vídeo

## **Fluxo Geral da Aplicação**
1. **Recebimento do Vídeo**: O usuário faz upload de um vídeo via API HTTP.
2. **Armazenamento Temporário**: O vídeo é salvo temporariamente no servidor.
3. **Registro no Banco de Dados**: A API registra a entrada do vídeo e seu status inicial (`pendente`).
4. **Processamento Concorrente**: Um worker em Go processa a conversão para HLS usando `ffmpeg`.
5. **Upload para AWS S3**: Os segmentos `.ts` e o manifesto `.m3u8` são enviados para um bucket na AWS.
6. **Atualização do Status**: Após o upload, o banco de dados é atualizado com o status `concluído` e a URL do vídeo.
7. **Consulta de Status**: O usuário pode consultar via API se a conversão foi concluída e obter o link final.

## **Desafios Técnicos e Soluções**
### **1. Gerenciamento de Concorrência**
- Uso de **worker pools** para processar múltiplos vídeos simultaneamente.
- Uso de **channels** para comunicação eficiente entre processos de conversão e upload.

## **Nome do Módulo**
`github.com/devfullcycle/golangtechweek`
