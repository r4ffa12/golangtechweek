# Configuração do FFmpeg para Conversão de Vídeos

## **Visão Geral**
O projeto utilizará o **FFmpeg** para processar e converter vídeos para o formato **HLS (HTTP Live Streaming)**. O FFmpeg será instalado dentro do **contêiner da aplicação**, garantindo que a conversão seja executada diretamente no ambiente containerizado.

## **Instalação do FFmpeg no Dockerfile**
- O **FFmpeg** será instalado no contêiner da aplicação, baseado na imagem **Golang 1.24 Alpine**.
- A instalação ocorrerá diretamente no **Dockerfile**, garantindo que todos os serviços tenham o FFmpeg disponível sem necessidade de instalação manual.

## **Execução do FFmpeg dentro do Contêiner**
- A conversão de vídeos será feita **inteiramente dentro do container**, eliminando dependências externas no ambiente local.
- O binário do **FFmpeg** será chamado diretamente dentro do código Go para processar os vídeos.
