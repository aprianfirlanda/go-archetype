---
name: fiber-swaggo-api
description: Use when adding, updating, documenting, or debugging Fiber HTTP APIs in the go-archetype repo, including handlers, request DTOs, response DTOs, validation, route registration, auth middleware, request ID responses, and Swaggo docs.
---

# Fiber Swaggo API

Use this skill for HTTP endpoint work.

## Files And Packages

- Handlers: `internal/adapters/http/handler/<domain>`; task package example is `taskhandler`.
- Request DTOs: `internal/adapters/http/dto/request/<domain>`; task example is `taskreq`.
- Response DTOs: `internal/adapters/http/dto/response/<domain>`; task example is `taskresp`.
- Shared response helpers: `internal/adapters/http/dto/response`.
- Shared request parsing helpers: `internal/adapters/http/dto/request`, package `httpreq`.
- Routes: `internal/adapters/http/router/router.go`.
- Validation: `internal/adapters/http/validation`.
- Server setup and global middleware: `internal/adapters/http/server/fiber.go`.

## Handler Pattern

- Get logger and request ID at the top:

```go
log := httpctx.Get(c, h.log)
rid := httpctx.GetRequestID(c)
```

- Parse path params in the handler.
- Use `httpreq.ParseBody[T]` and `httpreq.ParseQuery[T]` for body/query DTO parsing, validation, and standard 400 responses.
- Map DTOs into application `command` or `query` types.
- Call input ports through the handler service field using `c.UserContext()` (never replace with `context.Background()` in request path).
- Return shared response wrappers and include `rid`.
- Let application errors bubble to the global Fiber error handler when possible.

## Auth

Use existing middleware from `internal/adapters/http/middleware`:

- API key: `middleware.AuthAPIKey`
- JWT secret: `middleware.AuthJWT`
- Keycloak: `middleware.AuthKeycloak`
- Any accepted auth method: `middleware.AnyAuth(...)`

Follow existing route examples in `router.go`.

`middleware.AuthKeycloak` returns `(fiber.Handler, error)` because OIDC provider initialization can fail. Handle the error in route/server startup.

## Swagger

When an API route, DTO, response shape, or auth requirement changes:

1. Add/update Swaggo annotations on the handler.
2. Keep `@Security` aligned with route middleware.
3. Regenerate docs:

```sh
swag init -g cmd/http.go -o internal/adapters/http/docs
```

4. Inspect generated files in `internal/adapters/http/docs`.

## Verification

- Run `gofmt` on changed Go files.
- Run focused tests when available.
- Run `go test ./...` for broader confidence.
