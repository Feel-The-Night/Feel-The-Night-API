# Feel-The-Night API

Backend REST do projeto **Feel The Night**, uma plataforma de conteúdo para a comunidade de jogos de luta: guias, páginas de personagens e eventos.

## Sumário

- [Descrição Geral](#descrição-geral)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Estrutura do Repositório](#estrutura-do-repositório)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Configuração da Env](#configuração-da-env)
- [Rodando a Aplicação](#rodando-a-aplicação)
- [Documentação Swagger](#documentação-swagger)
- [Autenticação e Permissões](#autenticação-e-permissões)
- [Endpoints](#endpoints)
- [Banco de Dados](#banco-de-dados)
- [Testes](#testes)
- [Observações](#observações)
- [Licença](#licença)

## Descrição Geral

A API expõe endpoints para gerenciar usuários, personagens, guias e eventos, com autenticação via JWT e controle de acesso por papéis (iniciante, veterano e admin).

Ela é escrita em Go com Gin (HTTP), GORM (PostgreSQL) e Swaggo (documentação Swagger gerada a partir dos comentários dos handlers). O código segue uma organização **por feature**: cada recurso (usuário, personagem, guia, evento) tem a própria pasta com handler, service, DTOs e rotas, e a infraestrutura compartilhada fica em `internal/platform`.

## Tecnologias Utilizadas

- Go
- Gin
- GORM + PostgreSQL
- JWT (`golang-jwt/jwt`)
- bcrypt (`golang.org/x/crypto`)
- Rate limiting (`golang.org/x/time/rate`)
- Swaggo / Swag
- zap (logs estruturados)
- godotenv
- Air (hot reload em desenvolvimento)
- testify (testes)

## Estrutura do Repositório

```text
Feel-The-Night-Swagger/
├── cmd/
│   └── api/
│       └── main.go              # ponto de entrada: env, logger, banco e servidor HTTP
├── docs/                        # Swagger gerado pelo swag (não editar à mão)
├── internal/
│   ├── app/
│   │   └── routes.go            # NewRouter: monta middlewares, services, handlers e rotas
│   ├── auth/                    # login (handler, DTOs e rotas)
│   ├── user/                    # usuários
│   ├── character/               # personagens
│   ├── guide/                   # guias
│   ├── event/                   # eventos
│   │   ├── dto.go               #   entradas e respostas da feature
│   │   ├── handler.go           #   handlers HTTP
│   │   ├── routes.go            #   RegisterRoutes da feature
│   │   └── service.go           #   regras de negócio e acesso ao banco
│   ├── models/
│   │   └── models.go            # structs do GORM (compartilhadas entre as features)
│   └── platform/                # infraestrutura, sem regra de negócio
│       ├── apperr/              # erros compartilhados (ex.: ErrNoUpdateData)
│       ├── database/            # conexão, AutoMigrate, FKs e EscapeLike
│       ├── httpx/               # ErrorResponse, PageResponse e helpers de resposta/log
│       ├── jwt/                 # geração e validação de tokens
│       └── middleware/          # autenticação, roles, CORS, timeout e rate limit
├── tests/
│   └── seed_models.py           # script que popula a API com dados de exemplo
├── .air.toml
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

As dependências seguem uma direção só: `cmd/api → app → features → platform/models`. As features não importam umas às outras.

## Pré-requisitos

- Go (o `go.mod` está configurado com `go 1.27.1`)
- PostgreSQL em execução
- Git
- Opcional: [Air](https://github.com/air-verse/air) e [Swag](https://github.com/swaggo/swag)
- Opcional: Python 3, para rodar o script de seed

## Instalação

Na raiz do projeto:

```bash
go mod download
```

Crie o banco local:

```bash
createdb feel_the_night
```

Ferramentas opcionais:

```bash
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Configuração da Env

Copie o `.env.example` para `.env` e ajuste os valores:

```env
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=feel_the_night
DB_SSLMODE=disable

CORS_ORIGINS=http://localhost:3000,http://localhost:5173
JWT_SECRET=troque-por-um-segredo-aleatorio-com-pelo-menos-32-caracteres
```

| Variável | Obrigatória | Descrição |
|---|---|---|
| `JWT_SECRET` | sim | Segredo usado para assinar os tokens. Precisa ter **pelo menos 32 caracteres**, senão a API não sobe. Para gerar: `openssl rand -base64 48`. |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE` | sim | Conexão com o PostgreSQL. |
| `APP_PORT` | não | Porta local da API (padrão `8080`). Em produção, a variável `PORT` tem prioridade. |
| `APP_ENV` | não | `production` ativa a configuração de proxies confiáveis; qualquer outro valor é tratado como desenvolvimento. |
| `CORS_ORIGINS` | não | Origens liberadas no CORS, separadas por vírgula. Sem ela, valem as portas locais padrão do front. |
| `RENDER_EXTERNAL_HOSTNAME` | não | Definida automaticamente pelo Render; ajusta o host e o esquema (`https`) do Swagger em produção. |

O `.env` não é versionado. Em produção, as variáveis vêm do ambiente.

## Rodando a Aplicação

```bash
go run ./cmd/api
```

Com hot reload (regenera o Swagger e recompila a cada alteração):

```bash
air
```

Build:

```bash
go build -o ./bin/api ./cmd/api
```

## Documentação Swagger

Com a API rodando:

```text
http://localhost:8080/swagger/index.html
```

Para regenerar a documentação depois de alterar annotations ou DTOs:

```bash
swag init -d ./cmd/api,./internal/auth,./internal/user,./internal/character,./internal/guide,./internal/event,./internal/models,./internal/platform/httpx -g main.go
```

Rotas protegidas usam o esquema `BearerAuth`. No botão **Authorize**, cole o `access_token` recebido no login (com ou sem o prefixo `Bearer `).

## Autenticação e Permissões

### Fluxo

1. Crie uma conta em `POST /users`. Toda conta nova começa como **iniciante**.
2. Faça login em `POST /auth/login` com email e senha.
3. A resposta traz um `access_token` (válido por 15 minutos) e um `refresh_token` (válido por 7 dias).
4. Envie o access token nas rotas protegidas: `Authorization: Bearer <access_token>`.

Login e cadastro têm **rate limit por IP**: 5 tentativas seguidas e depois 1 a cada 12 segundos. Acima disso, a API responde `429` com o header `Retry-After`.

### Papéis

| Ação | Iniciante | Veterano | Admin |
|---|---|---|---|
| Ler qualquer conteúdo | ✅ | ✅ | ✅ |
| Criar guia | ❌ | ✅ | ✅ |
| Editar guia | ❌ | só os próprios | qualquer |
| Criar/editar personagem | ❌ | ✅ | ✅ |
| Criar/editar evento | ❌ | ✅ | ✅ |
| Editar a própria conta | ✅ | ✅ | ✅ (qualquer conta) |
| Deletar qualquer conteúdo | ❌ | ❌ | ✅ |
| Buscar usuário por email | ❌ | ❌ | ✅ |

### Primeiro admin

O papel nunca vem do corpo da requisição. Para criar o primeiro admin, cadastre a conta normalmente e promova direto no banco:

```sql
UPDATE users SET role = 'admin' WHERE email = 'seu@email.com';
```

Depois, faça login de novo: o papel vai dentro do token, então o token antigo continua com o papel anterior até expirar.

## Endpoints

| Método | Rota | Acesso | Descrição |
|---|---|---|---|
| `GET` | `/health` | público | Health check |
| `POST` | `/auth/login` | público (rate limit) | Login; retorna access e refresh token |
| `POST` | `/users` | público (rate limit) | Cadastro |
| `GET` | `/users?email=` | admin | Busca usuário por email |
| `GET` | `/users/:id` | público | Perfil público (sem email) |
| `PATCH` | `/users/:id` | o próprio usuário ou admin | Atualiza a conta |
| `DELETE` | `/users/:id` | admin | Remove usuário |
| `GET` | `/characters?name=&type=` | público | Lista personagens (filtros opcionais) |
| `GET` | `/characters/:id` | público | Busca personagem |
| `POST` | `/characters` | veterano ou admin | Cria personagem |
| `PATCH` | `/characters/:id` | veterano ou admin | Atualiza personagem |
| `DELETE` | `/characters/:id` | admin | Remove personagem (recusado com `409` se houver guias vinculados) |
| `GET` | `/guides?title=&page=&limit=` | público | Lista guias paginados |
| `GET` | `/guides/:id` | público | Busca guia |
| `POST` | `/guides` | veterano ou admin | Cria guia (o autor vem do token) |
| `PATCH` | `/guides/:id` | autor ou admin | Atualiza guia |
| `DELETE` | `/guides/:id` | admin | Remove guia |
| `GET` | `/events?title=&upcoming=&page=&limit=` | público | Lista eventos paginados (`upcoming=true` mostra só os futuros) |
| `GET` | `/events/:id` | público | Busca evento |
| `POST` | `/events` | veterano ou admin | Cria evento (o autor vem do token) |
| `PATCH` | `/events/:id` | veterano ou admin | Atualiza evento |
| `DELETE` | `/events/:id` | admin | Remove evento |

Detalhes de cada rota (corpo, respostas e códigos de erro) estão no Swagger.

### Formatos de resposta

Erros seguem sempre o mesmo formato:

```json
{ "error": "mensagem" }
```

Listagens paginadas (guias e eventos):

```json
{ "data": [], "page": 1, "limit": 20, "total": 0 }
```

Datas de evento usam o formato `AAAA-MM-DD`.

## Banco de Dados

A conexão é feita na subida da aplicação, com as variáveis `DB_*`. Ao iniciar, a API:

1. roda o `AutoMigrate` do GORM para todos os models de `internal/models`;
2. cria as foreign keys entre as tabelas (guias → usuários e personagens, eventos → usuários), caso ainda não existam.

As exclusões são **soft delete**: o registro recebe uma data em `deleted_at` e deixa de aparecer nas consultas, mas continua na tabela.

O `AutoMigrate` é prático em desenvolvimento, mas nunca remove colunas. Para produção com dados reais, o recomendado é migrar para migrations versionadas.

## Testes

Os testes de unidade ficam ao lado do código que testam (arquivos `*_test.go`) e não precisam de banco nem de Docker. Hoje cobrem:

- `internal/platform/jwt`: geração e validação de tokens (inclusive que um refresh token não passa como access token, e vice-versa);
- `internal/platform/middleware`: autenticação e checagem de papéis;
- `internal/platform/database`: escape de texto para buscas com `ILIKE`;
- `internal/event`: conversão de datas no formato `AAAA-MM-DD`.

### Rodar

```bash
go test ./...              # todos os testes
go test ./... -v           # mostra o nome e o resultado de cada teste
go test ./... -race        # detecta acesso concorrente sem lock (ex.: rate limiter)
go test ./... -cover       # mostra a cobertura por pacote
go test ./internal/platform/jwt/   # só um pacote
```

### Durante o desenvolvimento

Os testes **não** rodam no `air`, de propósito: o air recompila a cada save e esperar os testes deixaria o reload lento. Para acompanhar os testes enquanto programa, use o `gotestsum` num terminal separado:

```bash
go install gotest.tools/gotestsum@latest
gotestsum --watch
```

Ele roda novamente os testes do pacote que você alterou, sem interferir na API.

### Antes de enviar código

Rode `go test ./...` antes de commitar. Para não esquecer, dá para criar um hook de git que roda os testes antes de todo `git push` e cancela o push se algum falhar. Crie o arquivo `.git/hooks/pre-push`:

```bash
#!/bin/sh
go test ./... || exit 1
```

No Linux/macOS, dê permissão de execução com `chmod +x .git/hooks/pre-push`. Hooks ficam só na sua máquina (não são versionados), então cada pessoa do time precisa criar o seu.

### Dados de exemplo

Para popular uma API local com dados de exemplo:

```bash
pip install -r requirements.txt
python tests/seed_models.py
```

O script usa `BASE_URL` (padrão `http://localhost:8080`). Como as rotas de criação exigem veterano ou admin, a conta usada pelo seed precisa ter sido promovida antes.

## Observações

- A documentação em `docs/` é gerada pelo `swag`. Não edite esses arquivos à mão.
- Mensagens de erro internas (banco, timeout) são logadas com zap e nunca expostas ao cliente; ele recebe apenas `"erro interno"`.
- Toda requisição tem timeout de 10 segundos.
- Em produção (`APP_ENV=production`), as faixas de IP do proxy precisam estar configuradas em `internal/app/routes.go` para o rate limit identificar o IP real de cada cliente.

## Licença

Distribuído sob a licença descrita no arquivo [LICENSE](LICENSE).
