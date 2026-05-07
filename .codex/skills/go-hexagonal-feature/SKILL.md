---
name: go-hexagonal-feature
description: Use when adding or changing a domain feature in the go-archetype Go backend, especially tasks that span domain, application services, command/query/result models, input/output ports, HTTP adapters, GORM adapters, RabbitMQ adapters, Cobra wiring, migrations, tests, or Swagger docs.
---

# Go Hexagonal Feature

Use this skill for end-to-end feature work in the Go hexagonal architecture repo.

## First Checks

- Read root `AGENTS.md` first if present.
- Inspect package declarations before importing; package names often differ from folder names.
- Identify the target domain and whether the change is HTTP-only, messaging-only, persistence-only, or cross-layer.

## Layer Order

Prefer this order when implementing a feature:

1. Domain: add core entities, value objects, statuses, and domain errors under `internal/domain/<domain>`.
2. Application: add `command`, `query`, `result`, and `service` types under `internal/application/<domain>`.
3. Input ports: expose use case behavior from `internal/ports/input`.
4. Output ports: define persistence, publisher, or external service contracts in `internal/ports/output`.
5. Persistence adapter: implement GORM repositories under `internal/adapters/persistence/gorm/<domain>`.
6. HTTP adapter: add request DTOs, response DTOs, handlers, validation, and routes.
7. Messaging adapter: add RabbitMQ handlers/publishers when the feature is async.
8. Wiring: compose dependencies in `cmd/http.go`, `cmd/consumer.go`, or another owning Cobra command.
9. Migrations: add Goose migrations for schema changes.
10. Verification: run `gofmt`, focused tests, `go test ./...`, and `swag init` when API docs changed.

## Boundary Rules

- Domain must not import application, ports, adapters, GORM, Fiber, RabbitMQ, Cobra, Viper, or Logrus.
- Application services depend on ports and domain types.
- Handlers map DTOs to application commands/queries and map results to response DTOs.
- Repositories map between GORM models and domain entities inside the adapter.
- Wiring belongs in `cmd` and `internal/bootstrap`, not in domain/application packages.

## Context And Observability

- Pass `context.Context` through service, repository, publisher, and consumer calls.
- In HTTP paths, use `c.UserContext()` and keep request ID in responses.
- Use request-scoped loggers from `httpctx.Get(c, h.log)` in handlers.
- In services and port-out adapters, derive logger from context:

```go
log := logging.ComponentLogger(logging.FromContext(ctx), "<component>")
```

- Do not pass `rid` manually across service/repository/publisher calls; propagate by context.
- Use `db.WithContext(ctx)` in GORM operations.
- Verify `rid` continuity in both paths:
  - HTTP: middleware -> handler (`c.UserContext`) -> service -> repo/publisher
  - Consumer: message context (`CorrelationId`/`MessageId`) -> handler -> service -> repo/publisher

## Package Names To Remember

- HTTP context: `httpctx`
- Task request DTO: `taskreq`
- Task response DTO: `taskresp`
- HTTP task handler: `taskhandler`
- RabbitMQ root adapter: `messagingrmq`
- RabbitMQ task handler: `taskhandlermq`
- GORM root adapter: `gormadapter`
- GORM task repository: `taskgorm`
- Application command/query/result/service: `taskcmd`, `taskquery`, `taskresult`, `tasksvc`
- Ports: `portin`, `portout`
