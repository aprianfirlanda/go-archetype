# Go Archetype

Go Archetype is a backend service template built with hexagonal architecture. It keeps business logic in domain and application packages, with Fiber, GORM, RabbitMQ, Cobra, Viper, and Logrus isolated in adapters and infrastructure.

## Stack

- Go 1.25.4
- Cobra for CLI commands
- Viper for configuration
- Fiber for HTTP APIs
- GORM with PostgreSQL
- Goose for SQL migrations
- RabbitMQ via `github.com/rabbitmq/amqp091-go`
- Logrus for structured logging
- Swaggo for Swagger/OpenAPI docs

## Commands

Run the HTTP API:

```sh
go run . http
```

Run the RabbitMQ consumer:

```sh
go run . consumer
```

Run migrations:

```sh
go run . migrate up
go run . migrate status
go run . migrate version
go run . migrate down
```

Create a migration:

```sh
go run . migrate create add_example_table
```

## Local Setup

Start local dependencies:

```sh
docker compose up -d
```

Run the API:

```sh
go run . http --http-port 3000
```

Swagger UI is available at:

```text
http://localhost:3000/swagger/index.html
```

Manual HTTP request examples live in:

```text
test/http/request
```

## Configuration

Configuration is loaded through Cobra and Viper into `internal/infrastructure/config.Config`.

Precedence, highest to lowest:

1. Cobra flags
2. Environment variables
3. YAML config file
4. Flag defaults

Default config files:

- `config/config.yaml`
- `config/example.config.yaml`

Use a custom config file:

```sh
go run . http --config ./config/config.yaml
```

Environment variables use the app name as a prefix. With `appName = "go-archetype"`, the prefix is `GOARCHETYPE_`.

Examples:

```sh
GOARCHETYPE_HTTP_PORT=3000 go run . http
GOARCHETYPE_DB_HOST=localhost go run . migrate status
GOARCHETYPE_MESSAGING_RABBITMQ_URL=amqp://guest:guest@localhost:5672/ go run . consumer
```

When adding config:

- Add fields in `internal/infrastructure/config/schema.go`.
- Add flags in `cmd/root.go` for shared config or the owning command file for command-specific config.
- Keep flag names aligned with config keys: `http-port` maps to `http.port`.
- Update `config/example.config.yaml` when the value should be visible to users.
- Use the loaded `*config.Config`; avoid direct Viper reads outside config bootstrap and command setup.

## Architecture

The project follows ports and adapters:

- `internal/domain`: pure business concepts and rules.
- `internal/application`: use cases, commands, queries, results, and services.
- `internal/ports/input`: inbound use case interfaces.
- `internal/ports/output`: outbound contracts such as repositories, unit of work, publishers, and pingers.
- `internal/adapters/http`: Fiber handlers, DTOs, middleware, routing, validation, and Swagger docs.
- `internal/adapters/messaging/rabbitmq`: RabbitMQ connection, publisher, consumer, and message handlers.
- `internal/adapters/persistence/gorm`: GORM repositories, models, unit of work, pinger, and Goose migration adapter.
- `internal/infrastructure`: config, database bootstrap, logging, and test utilities.
- `cmd`: Cobra command wiring.

Keep dependencies pointing inward:

```text
adapters -> application -> domain
cmd/bootstrap -> adapters + application + infrastructure
domain -> no framework dependencies
application -> ports, not adapters
```

## HTTP API Conventions

HTTP code lives under `internal/adapters/http`.

Common flow:

1. Handler parses params/body/query.
2. Handler validates request DTOs.
3. Handler maps DTOs into application commands or queries.
4. Application service executes the use case through ports.
5. Handler maps results into response DTOs.

Use shared helpers for request parsing:

```go
req, err := httpreq.ParseBody[taskreq.Create](c, log, rid)
q, err := httpreq.ParseQuery[taskreq.List](c, log, rid)
```

Use shared response wrappers:

- `response.OK`
- `response.OKPaginate`
- `response.OKMessage`
- `response.Fail`
- `response.FailMessage`

Every HTTP response should include the request ID.

## Authentication

Supported middleware:

- API key: `Authorization: ApiKey <api_key>`
- JWT: `Authorization: Bearer <token>`
- Keycloak JWT: `Authorization: Bearer <token>`
- Combined auth with `middleware.AnyAuth(...)`

Keycloak middleware initializes an OIDC provider during route setup and returns an error if initialization fails.

## Messaging

RabbitMQ code lives under `internal/adapters/messaging/rabbitmq`.

Consumer retry behavior:

- Main queue: `<topic>`
- Retry queue: `<topic>.retry`
- DLQ: `<topic>.dlq`
- Retry count header: `x-retry-count`
- Retry error header: `x-last-error`
- DLQ headers: `x-final-error`, `x-original-queue`

Message handlers should unmarshal, map to application commands, call input ports, and return errors so the consumer can ack/nack correctly.

## Persistence

GORM persistence code lives under `internal/adapters/persistence/gorm`.

Conventions:

- Repository interfaces live in `internal/ports/output`.
- Repositories accept `context.Context`.
- GORM calls use `db.WithContext(ctx)`.
- GORM models and domain mapping stay inside the persistence adapter.
- Schema changes use Goose migrations in `migrations/`.

## Swagger

Regenerate Swagger docs after API, DTO, response, or auth annotation changes:

```sh
swag init -g cmd/http.go -o internal/adapters/http/docs
```

Generated files:

- `internal/adapters/http/docs/docs.go`
- `internal/adapters/http/docs/swagger.json`
- `internal/adapters/http/docs/swagger.yaml`

## Verification

Format changed Go files:

```sh
gofmt -w <changed-go-files>
```

Run the package suite:

```sh
go test ./...
```

At the moment this repository does not include committed `_test.go` files, so `go test ./...` is primarily a compile check until tests are added.
