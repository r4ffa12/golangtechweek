# Configuração do Ambiente com Docker e Docker Compose

## **Visão Geral**
Este projeto será executado inteiramente dentro de contêineres Docker, eliminando a necessidade de instalar **Go** e **PostgreSQL** localmente. Para isso, utilizaremos **Docker Compose** para orquestrar os serviços da aplicação e do banco de dados.

## **Tecnologias Utilizadas**
- **Docker** para contêinerização do ambiente.
- **Docker Compose** para gerenciar múltiplos serviços.
- **Imagem Base da Aplicação:** `golang:1.24-alpine` para otimização e leveza.
- **Banco de Dados:** PostgreSQL 14 rodando em contêiner.

## **Estrutura dos Serviços no Docker Compose**
- **`app`**: Serviço principal onde a aplicação Go será executada dentro do contêiner.
- **`db`**: Serviço PostgreSQL para armazenamento das informações de conversão.
- **`localstack`**: Serviço Localstack para testar o upload para S3 localmente.

## **Execução Totalmente Dentro do Container**
- O contêiner da aplicação terá o ambiente Go configurado para rodar diretamente dentro dele.
- O código será montado como um volume para desenvolvimento dinâmico sem necessidade de reconstrução manual do container.
- O serviço será executado com **TTY habilitado** (`tty: true`) para manter o processo ativo.
- O serviço **`localstack`** será executado com **TTY habilitado** (`tty: true`) para manter o processo ativo.

## **Bind Mounts**
- O código da aplicação será montado como um bind mount para desenvolvimento, logo, qualquer alteração no código será refletida automaticamente no container.
- O banco de dados será montado como um volume para persistência dos dados.

## **Arquivo `docker-compose.yaml` (Descrição Geral)**
- Define os serviços **`app`** e **`db`** e **`localstack`**.
- Configura um volume para persistência dos dados do PostgreSQL.
- Exibe logs da aplicação e do banco em tempo real.