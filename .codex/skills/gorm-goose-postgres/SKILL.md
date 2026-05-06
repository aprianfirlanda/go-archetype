---
name: gorm-goose-postgres
description: Use when adding, changing, or debugging PostgreSQL persistence in the go-archetype repo, including output repository ports, GORM models and mappers, repository implementations, unit of work usage, Goose migrations, DB config, and persistence tests.
---

# GORM Goose Postgres

Use this skill for database persistence and migration work.

## Files And Packages

- Output ports: `internal/ports/output`, package `portout`.
- GORM root adapter: `internal/adapters/persistence/gorm`, package `gormadapter`.
- Domain repositories: `internal/adapters/persistence/gorm/<domain>`.
- Task repository example package: `taskgorm`.
- Migrations: `migrations/`.
- Goose wrapper: `internal/adapters/persistence/gorm/migrate`.
- DB bootstrap/config: `internal/infrastructure/db` and `internal/infrastructure/config`.

## Repository Pattern

- Define contracts in output ports before adapter code.
- Keep GORM models, table names, indexes, and domain mapping inside the GORM adapter.
- Repository methods must accept `context.Context`.
- Use `r.db.WithContext(ctx)` for queries.
- Return domain entities from repositories, not GORM models.
- Wrap repository errors at application service level with `apperror` when appropriate.

## Unit Of Work

- Use existing `portout.UnitOfWork` and `gormadapter.NewUnitOfWork` for transactional use cases.
- Keep transaction concerns in application service orchestration or adapter code, matching existing patterns.

## Migrations

Create migrations with the Cobra command:

```sh
go run . migrate create <name>
```

Other useful commands:

```sh
go run . migrate up
go run . migrate down
go run . migrate status
go run . migrate version
```

When schema changes:

- Add a Goose migration under `migrations/`.
- Keep SQL reversible when possible.
- Update GORM model fields and mapping functions.
- Add/update repository tests for behavior, not just model shape.

## Verification

- Run `gofmt` on changed Go files.
- Run focused repository tests first.
- Run `go test ./...` when database/test environment is available.

