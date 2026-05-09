# Project Structure

This repository uses hexagonal architecture. The main rule is simple: domain and application code stay independent from frameworks, while adapters translate between external systems and use cases.

## Top Level

```text
.
|-- cmd/                  Cobra commands and dependency wiring
|-- config/               YAML configuration examples/defaults
|-- internal/             Private application code
|-- migrations/           Goose SQL migrations
|-- test/http/request/    Manual HTTP request examples
|-- compose.yaml          Local dependency stack
|-- Dockerfile            Container build
|-- main.go               CLI entrypoint
|-- README.md             Project overview and usage
`-- PROJECT_STRUCTURE.md  This file
```

## Command Layer

```text
cmd/
|-- root.go       Root command, persistent flags, config/logger/DB/RabbitMQ bootstrap
|-- http.go       HTTP server command and HTTP dependency wiring
|-- consumer.go   RabbitMQ consumer command and consumer dependency wiring
`-- migrate.go    Goose migration command wiring
```

`cmd` is allowed to import infrastructure, adapters, and application packages because it composes the application. Do not move wiring into domain or application packages.

## Domain Layer

```text
internal/domain/
|-- auth/       Auth-related domain types
|-- identity/   Public ID generation and identity concepts
`-- task/       Task entity, status, and domain errors
```

Rules:

- Pure Go business concepts only.
- No Fiber, GORM, RabbitMQ, Cobra, Viper, or Logrus imports.
- Domain errors and invariants belong here.

## Application Layer

```text
internal/application/
|-- health/
|   `-- service/
`-- task/
    |-- command/   Write-use-case inputs
    |-- query/     Read-use-case filters
    |-- result/    Use-case output models
    `-- service/   Use-case orchestration
```

Rules:

- Application services depend on domain types and ports.
- Services should not import HTTP, GORM, RabbitMQ, Cobra, Viper, or other adapters.
- Use `context.Context` throughout service methods.

## Ports

```text
internal/ports/
|-- input/
|   |-- health_service.go
|   |-- message_consumer.go
|   `-- task_service.go
`-- output/
    |-- db_pinger.go
    |-- message_publisher.go
    |-- repository.go
    `-- uow.go
```

Input ports describe use cases exposed to inbound adapters. Output ports describe capabilities required by application services, such as persistence, unit of work, health checks, and publishing.

## HTTP Adapter

```text
internal/adapters/http/
|-- context/       Fiber context helpers for request ID, logger, and user data
|-- docs/          Generated Swagger/OpenAPI files
|-- dto/
|   |-- request/   Request DTOs by domain, plus shared body/query parsing helpers
|   `-- response/  Response DTOs and shared wrappers
|-- handler/
|   |-- demo/      Demo endpoints
|   `-- task/      Task endpoints
|-- middleware/    Auth, request ID, logging, recovery, and health middleware
|-- router/        Route registration
|-- server/        Fiber app setup
`-- validation/    Validator setup and field error messages
```

Handler responsibilities:

- Read request ID and request-scoped logger.
- Parse params/body/query.
- Validate DTOs.
- Map to application commands or queries.
- Call input ports with `c.UserContext()`.
- Return response DTOs using shared response helpers.

Shared request helpers:

```go
httpreq.ParseBody[T](c, log, rid)
httpreq.ParseQuery[T](c, log, rid)
```

## Messaging Adapter

```text
internal/adapters/messaging/rabbitmq/
|-- bootstrap.go
|-- connection.go
|-- consumer.go
|-- publisher.go
`-- handler/
    `-- task/
```

Responsibilities:

- Manage RabbitMQ connection/channel setup.
- Publish messages through output ports.
- Consume topics through input consumer registration.
- Map message payloads into application commands.
- Preserve request/correlation ID through AMQP metadata.

## Persistence Adapter

```text
internal/adapters/persistence/gorm/
|-- migrate/       Goose migrator wrapper
|-- task/          Task GORM model, mapper, and repository methods
|-- pinger.go      DB health pinger
|-- uow.go         Unit of work implementation
`-- uow_tx.go      Transaction implementation
```

Rules:

- GORM models do not leave this adapter.
- Repositories return domain entities.
- Repository methods accept `context.Context`.
- Use `db.WithContext(ctx)` for queries.
- Application services wrap adapter errors into application errors where needed.

## Infrastructure

```text
internal/infrastructure/
|-- config/     Config schema, Viper setup, and config loading
|-- db/         PostgreSQL connection, pool setup, ping, and GORM logger
|-- logging/    Logrus setup and context-aware logger helpers
`-- testutil/   Shared test utilities
```

Infrastructure contains cross-cutting implementation details. Keep direct Viper usage inside `config` and command bootstrap.

## Bootstrap

```text
internal/bootstrap/
|-- consumer_registry.go
`-- http_app.go
```

Bootstrap package types group dependencies for adapters and commands. It should stay lightweight and avoid business logic.

## Shared Packages

```text
internal/pkg/
|-- apperror/   Application error type and helpers
`-- jwtutil/    JWT helper functions
```

Use `internal/pkg` for small shared helpers that are not domain concepts and do not belong to a specific adapter.

## Adding A Domain Feature

Use this order for new features:

1. Add or update domain types under `internal/domain/<domain>`.
2. Add command/query/result models under `internal/application/<domain>`.
3. Implement use cases under `internal/application/<domain>/service`.
4. Add input/output ports under `internal/ports`.
5. Implement persistence under `internal/adapters/persistence/gorm/<domain>` if needed.
6. Add request/response DTOs under `internal/adapters/http/dto`.
7. Add handlers under `internal/adapters/http/handler/<domain>`.
8. Register routes in `internal/adapters/http/router/router.go`.
9. Wire dependencies in the owning Cobra command.
10. Add Goose migrations for schema changes.
11. Regenerate Swagger docs for API changes.

## Useful Inspection Commands

List repository files:

```sh
rg --files
```

Show Go files:

```sh
rg --files -g '*.go'
```

Show package declarations:

```sh
rg '^package ' internal cmd
```

Show imports of a framework:

```sh
rg 'github.com/gofiber|gorm.io/gorm|github.com/spf13/viper|github.com/sirupsen/logrus' internal cmd
```

Display directory tree when `tree` is installed:

```sh
tree -I '.git|bin|tmp' -L 4
```
