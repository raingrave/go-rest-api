# API REST com Go, Gin e PostgreSQL

Uma API RESTful simples e robusta construída com Go, o framework Gin e PostgreSQL. Este projeto é totalmente containerizado usando Docker e Docker Compose, proporcionando um ambiente de desenvolvimento limpo e reproduzível.

## ✨ Funcionalidades

- **Autenticação JWT:** Endpoints seguros com JSON Web Tokens.
- **Gestão de Usuários:** Funcionalidade CRUD completa para usuários, com senhas criptografadas e validação de entrada.
- **Migrações Automáticas de BD:** O schema do banco de dados é gerenciado e aplicado automaticamente ao iniciar a aplicação.
- **Health Check:** Um endpoint `/health` para monitorar o status da API.
- **Suíte de Testes:** Cobertura de testes para garantir a confiabilidade dos handlers e middlewares.
- **Containerização:** Roda inteiramente em containers Docker para consistência e facilidade de implantação.
- **Estrutura Organizada:** Segue o layout de projeto padrão do Go.

## 🛠️ Tecnologias Utilizadas

- **Go:** Linguagem de programação principal.
- **Gin:** Framework web HTTP de alta performance.
- **PostgreSQL:** Sistema de banco de dados objeto-relacional.
- **sqlx:** Extensões para o pacote `database/sql` padrão.
- **golang-jwt:** Para geração e validação de tokens JWT.
- **bcrypt:** Para hashing seguro de senhas.
- **golang-migrate:** Para o gerenciamento de migrações do banco de dados.
- **testify:** Biblioteca de asserções para testes.
- **Docker & Docker Compose:** Para containerizar e orquestrar os serviços.

## 🚀 Começando

Siga estas instruções para obter uma cópia do projeto e executá-lo em sua máquina local.

### Pré-requisitos

- [Go](https://go.dev/doc/install) (v1.24 ou mais recente)
- [Docker](https://docs.docker.com/get-docker/) e [Docker Compose](https://docs.docker.com/compose/install/)

### Instalação e Execução

1.  **Clone o repositório:**
    ```sh
    git clone git@github.com:raingrave/go-rest-api.git
    cd go-rest-api
    ```

2.  **Crie o arquivo de ambiente:**
    Copie o arquivo de exemplo `.env.example` para um novo arquivo chamado `.env`.
    ```sh
    cp .env.example .env
    ```
    *Você pode ajustar os valores no arquivo `.env` se necessário.*

3.  **Execute a aplicação com Docker Compose:**
    Este comando irá construir a imagem da API, iniciar os containers e conectá-los. Ao iniciar, a aplicação irá rodar as **migrações do banco de dados automaticamente**, configurando o schema necessário.
    ```sh
    docker compose up --build -d
    ```
    A API estará disponível em `http://localhost:3000`.

## 🧪 Executando os Testes

Para rodar a suíte de testes completa, certifique-se de que os containers do Docker estejam em execução (pois os testes precisam de uma conexão com o banco de dados) e execute o seguinte comando na raiz do projeto:

```sh
go test ./...
```

Para ver a cobertura de testes em detalhe:

```sh
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out
```

## Endpoints da API

A URL base é `http://localhost:3000/api/v1`.

| Método   | Endpoint      | Descrição                               | Autenticação | Corpo da Requisição (Exemplo)                     |
| :------- | :------------ | :-------------------------------------- | :----------- | :------------------------------------------------ |
| `GET`    | `/health`     | Verifica o status da API (fora da versão).| Nenhuma      | `N/A`                                             |
| `POST`   | `/users`      | Cria um novo usuário. Validações: `name` (obrigatório), `email` (obrigatório, formato de email), `password` (obrigatório, min 8 caracteres). | Nenhuma      | `{"name":"...", "email":"...", "password":"..."}` |
| `POST`   | `/login`      | Autentica um usuário e retorna um token. | Nenhuma      | `{"email":"...", "password":"..."}`               |
| `GET`    | `/users`      | Lista todos os usuários.                | **Bearer Token** | `N/A`                                             |
| `GET`    | `/users/{id}` | Busca um único usuário pelo ID.         | **Bearer Token** | `N/A`                                             |
| `PUT`    | `/users/{id}` | Atualiza um usuário existente.          | **Bearer Token** | `{"name":"...", "email":"..."}`                   |
| `DELETE` | `/users/{id}` | Deleta um usuário pelo ID.              | **Bearer Token** | `N/A`                                             |

### Como se Autenticar

1.  Crie um usuário via `POST /users`.
2.  Faça login com as credenciais via `POST /login` para receber um token. O tempo de expiração padrão do token é de 60 minutos, mas pode ser configurado através da variável de ambiente `JWT_EXPIRATION_MINUTES`.
3.  Para acessar os endpoints protegidos, inclua o cabeçalho `Authorization` em suas requisições:
    ```
    Authorization: Bearer <seu_token_jwt_aqui>
    ```

### Erros de Validação

Ao criar ou atualizar recursos, se houver um erro de validação nos dados enviados, a API retornará uma resposta `400 Bad Request` com um corpo JSON detalhando os campos problemáticos.

**Exemplo de Resposta de Erro:**
```json
{
    "errors": {
        "Email": "Invalid email format",
        "Password": "This field must be at least 8 characters long"
    }
}
```

## 🏛️ Arquitetura e Fluxo de Dados

Os diagramas de sequência abaixo ilustram os principais fluxos da aplicação.

### 1. Criação de Usuário

```mermaid
sequenceDiagram
    participant Client
    participant Gin Router
    participant User Handler
    participant Bcrypt
    participant User Repository
    participant PostgreSQL DB

    Client->>Gin Router: POST /api/v1/users (com nome, email, senha)
    Gin Router->>User Handler: Chama CreateUser(c)
    User Handler->>User Handler: Valida os dados de entrada (nome, email, senha)
    User Handler->>Bcrypt: GenerateFromPassword(senha em texto plano)
    Bcrypt-->>User Handler: Retorna senha com hash
    User Handler->>User Repository: CreateUser(usuário com senha hasheada)
    User Repository->>PostgreSQL DB: INSERT INTO users (...) VALUES (...)
    PostgreSQL DB-->>User Repository: Retorna o novo UUID do usuário
    User Repository-->>User Handler: Retorna o UUID
    User Handler-->>Gin Router: Resposta 201 Created (com dados do usuário, sem senha)
    Gin Router-->>Client: Resposta 201 Created
```

### 2. Autenticação (Login)

```mermaid
sequenceDiagram
    participant Client
    participant Gin Router
    participant Auth Handler
    participant User Repository
    participant PostgreSQL DB
    participant Bcrypt
    participant JWT Library

    Client->>Gin Router: POST /api/v1/login (com email e senha)
    Gin Router->>Auth Handler: Chama Login(c)
    Auth Handler->>Auth Handler: Valida os dados de entrada (email, senha)
    Auth Handler->>User Repository: GetUserByEmail(email)
    User Repository->>PostgreSQL DB: SELECT * FROM users WHERE email = ...
    PostgreSQL DB-->>User Repository: Retorna dados do usuário (incluindo senha com hash)
    User Repository-->>Auth Handler: Retorna o objeto User
    Auth Handler->>Bcrypt: CompareHashAndPassword(hash do DB, senha da requisição)
    Bcrypt-->>Auth Handler: Confirma que a senha é válida
    Auth Handler->>JWT Library: Gera o token com ID do usuário e tempo de expiração
    JWT Library-->>Auth Handler: Retorna o token assinado (string)
    Auth Handler-->>Gin Router: Resposta 200 OK (com o token)
    Gin Router-->>Client: Resposta 200 OK
```

### 3. Acesso a Recurso Protegido

```mermaid
sequenceDiagram
    participant Client
    participant Gin Router
    participant Auth Middleware
    participant JWT Library
    participant User Handler
    participant User Repository
    participant PostgreSQL DB

    Client->>Gin Router: GET /api/v1/users (com Header "Authorization: Bearer <token>")
    Gin Router->>Auth Middleware: Executa o middleware de autenticação
    Auth Middleware->>Auth Middleware: Extrai o token do cabeçalho
    Auth Middleware->>JWT Library: ParseAndValidate(token)
    JWT Library-->>Auth Middleware: Retorna que o token é válido (assinatura e expiração OK)
    Auth Middleware->>Gin Router: Chama c.Next() para continuar
    Gin Router->>User Handler: Chama ListUsers(c)
    User Handler->>User Repository: ListUsers()
    User Repository->>PostgreSQL DB: SELECT * FROM users
    PostgreSQL DB-->>User Repository: Retorna a lista de usuários
    User Repository-->>User Handler: Retorna a lista
    User Handler-->>Gin Router: Resposta 200 OK (com a lista de usuários)
    Gin Router-->>Client: Resposta 200 OK
```