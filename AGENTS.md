# AGENTS.md

Guidance for coding agents working in this repository.

## Project Overview

This is a Go backend archetype using hexagonal architecture with ports and adapters.

Primary stack:

- Go 1.25.4
- Cobra CLI for command management
- Viper for config loading
- Fiber for HTTP
- RabbitMQ via `github.com/rabbitmq/amqp091-go`
- GORM with PostgreSQL
- Goose migrations
- Logrus structured logging
- Swaggo for Swagger/OpenAPI docs

Available CLI commands:

- `go run . http`
- `go run . consumer`
- `go run . migrate <create|up|up-to|down|down-to|status|version>`

## Architecture Rules

Keep the hexagonal boundaries intact:

- `internal/domain/<domain>` contains core business concepts and rules. Keep it framework-free.
- `internal/application/<domain>` contains use cases split into `service`, `command`, `query`, and `result`.
- `internal/ports/input` contains inbound use case interfaces.
- `internal/ports/output` contains outbound interfaces such as repositories, UoW, DB pingers, and message publishers.
- `internal/adapters/http` contains Fiber-specific request handling, DTOs, validation, middleware, router, and server setup.
- `internal/adapters/messaging/rabbitmq` contains RabbitMQ connection, publisher, consumer, and messaging handlers.
- `internal/adapters/persistence/gorm` contains GORM/PostgreSQL implementations.
- `internal/bootstrap` contains dependency structs and registries used by commands.
- `cmd` is the Cobra command layer and is responsible for application wiring.

Do not import adapters into domain or application logic. Application services should depend on ports, not on Fiber, GORM, RabbitMQ, Cobra, Viper, or Logrus.

## Package Names

Always inspect package declarations before importing. Package names are sometimes intentionally different from folder names.

Common examples:

- `internal/adapters/http/context` uses package `httpctx`
- `internal/adapters/http/dto/request/task` uses package `taskreq`
- `internal/adapters/http/dto/response/task` uses package `taskresp`
- `internal/adapters/http/handler/task` uses package `taskhandler`
- `internal/adapters/http/handler/demo` uses package `demohandler`
- `internal/adapters/messaging/rabbitmq` uses package `messagingrmq`
- `internal/adapters/messaging/rabbitmq/handler/task` uses package `taskhandlermq`
- `internal/adapters/persistence/gorm` uses package `gormadapter`
- `internal/adapters/persistence/gorm/task` uses package `taskgorm`
- `internal/application/task/command` uses package `taskcmd`
- `internal/application/task/query` uses package `taskquery`
- `internal/application/task/result` uses package `taskresult`
- `internal/application/task/service` uses package `tasksvc`
- `internal/application/health/service` uses package `healthsvc`
- `internal/ports/input` uses package `portin`
- `internal/ports/output` uses package `portout`

## Configuration

Configuration is centralized in `internal/infrastructure/config`.

- Schema lives in `internal/infrastructure/config/schema.go`.
- Loading/unmarshal lives in `internal/infrastructure/config/config.go`.
- Cobra/Viper initialization lives in `internal/infrastructure/config/bootstrap.go`.
- Root and command flags live in `cmd/root.go` and command files such as `cmd/http.go`.

Viper merges configuration sources:

1. Cobra flags
2. Environment variables
3. YAML config file
4. Flag defaults

Environment variable prefix depends on the `appName` variable in `cmd/root.go`.
The config bootstrap removes `-` from `appName` and uppercases the result with `strings.ToUpper(strings.ReplaceAll(appName, "-", ""))`.
Examples:

- `appName = "go-archetype"` uses prefix `GOARCHETYPE_`, so `http.port` can be set with `GOARCHETYPE_HTTP_PORT`.
- `appName = "clm-be"` uses prefix `CLMBE_`, so `http.port` can be set with `CLMBE_HTTP_PORT`.

When adding config:

- Add fields to `config.Config` and nested structs with `mapstructure` tags.
- Add Cobra flags in the root command for shared config or the owning subcommand for command-specific config.
- Keep config keys aligned with Viper flag names, where dashes become dots.
- Update `config/example.config.yaml` and `config/config.yaml` if the new value must be visible in file config.
- Avoid reading Viper directly outside the config package and command bootstrap. Use the loaded `*config.Config`.

## Dependency Wiring

Wiring belongs in Cobra/root command flow, not inside domain or application packages.

- `cmd/root.go` initializes config, logger, DB, and RabbitMQ in `PersistentPreRunE`.
- `cmd/http.go` wires repositories, unit of work, application services, and `bootstrap.HttpApp`.
- `cmd/consumer.go` wires repositories, unit of work, services, message handlers, and the consumer registry.
- `cmd/migrate.go` wires Goose migration commands.

For new dependencies, add ports first, implement them in adapters, then wire implementations in the relevant command.

## HTTP APIs

HTTP adapter conventions:

- Handlers: `internal/adapters/http/handler/<domain>`
- Request DTOs: `internal/adapters/http/dto/request/<domain>`
- Response DTOs: `internal/adapters/http/dto/response/<domain>`
- Shared response wrappers: `internal/adapters/http/dto/response`
- Shared request parsing helpers: `internal/adapters/http/request`
- Routes: `internal/adapters/http/router/router.go`
- Validation: `internal/adapters/http/validation`
- Server/middleware setup: `internal/adapters/http/server/fiber.go`

When adding or changing an API:

- Keep Fiber-specific parsing and response formatting in handlers.
- Use `request.ParseBody[T]` and `request.ParseQuery[T]` to avoid repeating parse/validation/error response boilerplate.
- Convert request DTOs into application `command` or `query` types.
- Return response DTOs, not domain entities, unless an existing endpoint already does so.
- Use `response.OK`, `response.OKPaginate`, `response.OKMessage`, `response.Fail`, or `response.FailMessage`.
- Include the request ID in every response.
- Add or update Swaggo annotations on handlers.
- Regenerate Swagger docs after API changes:

```sh
swag init -g cmd/http.go -o internal/adapters/http/docs
```

## Authentication

HTTP routes can be protected with:

- Keycloak JWT: `middleware.AuthKeycloak`
- Local JWT secret: `middleware.AuthJWT`
- API key: `middleware.AuthAPIKey`
- Combined auth: `middleware.AnyAuth(...)`

Follow current route patterns in `internal/adapters/http/router/router.go`.

Current API key middleware expects:

```http
Authorization: ApiKey <api_key>
```

JWT and Keycloak middleware use bearer tokens:

```http
Authorization: Bearer <token>
```

`middleware.AuthKeycloak` initializes an OIDC provider and returns `(fiber.Handler, error)`. Handle that error during route/server startup instead of panicking inside middleware construction.

## Logging And Request ID

Request ID is mandatory for observability.

- Fiber request IDs are generated by `requestid.New()` and read through `httpctx.GetRequestID(c)`.
- Request-scoped loggers are stored through `httpctx.Set` and read with `httpctx.Get`.
- HTTP handlers should log with the request-scoped logger and include `rid` in responses.
- Application services and repositories receive `context.Context`; do not replace it with `context.Background()` in request paths.
- In service and port-out methods, derive logger from context with `logging.ComponentLogger(logging.FromContext(ctx), "<component>")`.
- Do not pass `rid` manually to service/repository/publisher logs; use context-based logger enrichment.
- GORM calls should use `db.WithContext(ctx)`.
- RabbitMQ publish/consume paths should pass context through the port methods.
- HTTP flow: request middleware puts request logger and `rid` into `c.UserContext()`, handlers call services with `c.UserContext()`, then service/repository/publisher logs use that context.
- Consumer flow: consumer creates per-message context using `CorrelationId` or `MessageId` as `rid`, handlers call services with that context, then service/repository/publisher logs use that context.
- Publisher should set AMQP `CorrelationId` and `MessageId` from context `rid`; consumer should read them to rehydrate context `rid`.
- New logs should use Logrus structured fields and include the existing request/correlation ID whenever available.
- Use `logging.ComponentLogger` or `logging.WithComponentAndFields` for stable component names.
- Keep both `rid` and `request_id` fields in logs for compatibility.

For HTTP handlers, the usual pattern is:

```go
log := httpctx.Get(c, h.log)
rid := httpctx.GetRequestID(c)
```

## Messaging

RabbitMQ adapter code lives in `internal/adapters/messaging/rabbitmq`.

- Publisher port: `internal/ports/output/message_publisher.go`
- Consumer port: `internal/ports/input/message_consumer.go`
- Messaging handlers: `internal/adapters/messaging/rabbitmq/handler/<domain>`
- Consumer registration: `internal/bootstrap/consumer_registry.go`
- Command wiring: `cmd/consumer.go`

Keep message handlers thin: unmarshal payloads, map them to application commands, call input ports, and return errors so the consumer can ack/nack correctly.

Consumer retry/backoff behavior:

- Main queue name: `<topic>`
- Retry queue name: `<topic>.retry`
- DLQ name: `<topic>.dlq`
- Retry queue should dead-letter back to main queue with:
  - `x-dead-letter-exchange: ""`
  - `x-dead-letter-routing-key: <topic>`
- Retry count is tracked in message header `x-retry-count`.
- Failure metadata should be tracked in headers:
  - `x-last-error` for retry publish
  - `x-final-error` and `x-original-queue` for DLQ publish
- Backoff delays are taken from `cfg.Messaging.RabbitMQ.Consumer.Retry.Backoff`.
- If `retry_count` exceeds configured backoff slots, message should be published to DLQ and original message acknowledged.
- On retry publish failure or DLQ publish failure, `Nack(false, true)` the original message.
- Preserve message metadata (`CorrelationId`, `MessageId`, headers, body, content type) when republishing.

## Persistence And Migrations

GORM persistence implementations live under `internal/adapters/persistence/gorm`.

- Domain-specific repositories go under `internal/adapters/persistence/gorm/<domain>`.
- Repository interfaces belong in `internal/ports/output`.
- Use `context.Context` in every repository method.
- Keep GORM models and domain mapping inside the persistence adapter.

Migrations live in `migrations/` and are managed by Goose through the Cobra migrate command.

Useful commands:

```sh
go run . migrate create <name>
go run . migrate up
go run . migrate down
go run . migrate status
go run . migrate version
```

## Adding A New Domain Feature

Typical flow:

1. Add or update domain entities/errors/value objects in `internal/domain/<domain>`.
2. Add application command/query/result types under `internal/application/<domain>`.
3. Add service behavior under `internal/application/<domain>/service`.
4. Add or update input ports in `internal/ports/input`.
5. Add or update output ports in `internal/ports/output`.
6. Implement persistence under `internal/adapters/persistence/gorm/<domain>`.
7. Add HTTP request/response DTOs under `internal/adapters/http/dto`.
8. Add HTTP handler methods under `internal/adapters/http/handler/<domain>`.
9. Register routes and auth middleware in `internal/adapters/http/router/router.go`.
10. Wire dependencies in the relevant Cobra command.
11. Add migrations when schema changes.
12. Regenerate Swagger docs for API changes.
13. Add focused tests for domain/application/persistence behavior.

## Testing And Verification

Use focused verification first:

```sh
go test ./...
```

For Swagger-impacting changes, run:

```sh
swag init -g cmd/http.go -o internal/adapters/http/docs
```

For formatting:

```sh
gofmt -w <changed-go-files>
```

Do not commit generated docs or migration changes blindly. Inspect diffs and keep changes scoped to the requested task.
