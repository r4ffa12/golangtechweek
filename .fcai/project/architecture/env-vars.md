# Uso de Variáveis de Ambiente no Projeto

## **Visão Geral**

O projeto utilizará **variáveis de ambiente** para tornar a configuração mais flexível, segura e adaptável a diferentes ambientes de execução. Em vez de definir valores sensíveis diretamente no código, utilizaremos variáveis para armazenar informações como **parâmetros do banco de dados**, **chaves de API** e **credenciais de serviços externos**.

## **Motivação para Uso de Variáveis de Ambiente**

- **Flexibilidade**: Permite configurar o sistema para diferentes ambientes (desenvolvimento, teste, produção) sem modificar o código.
- **Segurança**: Evita armazenar credenciais sensíveis no código-fonte, reduzindo riscos de exposição acidental.
- **Facilidade de Deploy**: As configurações podem ser alteradas diretamente via `.env`&#x20;

## **Principais Configurações Usadas no Projeto**

### **Exemplos de Banco de Dados**

- `DB_HOST`: Endereço do servidor PostgreSQL
- `DB_PORT`: Porta do banco de dados
- `DB_USER`: Usuário do banco de dados
- `DB_PASSWORD`: Senha do banco de dados
- `DB_NAME`: Nome do banco de dados

### **Exemplos de Serviços Externos e APIs**

- `AWS_ACCESS_KEY_ID`: Chave de acesso AWS
- `AWS_SECRET_ACCESS_KEY`: Chave secreta AWS
- `S3_BUCKET_NAME`: Nome do bucket S3
- `OPENAI_API_KEY`: Chave de API para uso do OpenAI
- `EMAIL_SERVICE_API_KEY`: Chave de API para serviços de e-mail (ex: SendGrid, Mailgun)

### **Configurações Gerais**

- `APP_ENV`: Define o ambiente de execução (`development`, `staging`, `production`)
- `LOG_LEVEL`: Define o nível de log da aplicação (`debug`, `info`, `warn`, `error`)

## **Carregamento das Variáveis no Ambiente Docker**

- As variáveis serão carregadas automaticamente a partir de um arquivo `.env`, garantindo que a configuração possa ser facilmente alterada sem modificar os arquivos do código.
- No **Docker Compose**, as variáveis serão passadas para os contêineres, os valores default podem ser especificados no serviço através da seção: environment 
