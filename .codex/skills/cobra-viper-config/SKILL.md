---
name: cobra-viper-config
description: Use when adding, changing, or debugging configuration in the go-archetype repo, including Cobra flags, Viper binding, YAML config files, environment variables, config schema structs, appName-derived env prefixes, and command pre-run config loading.
---

# Cobra Viper Config

Use this skill for configuration changes.

## Where Config Lives

- `cmd/root.go`: `appName`, root persistent flags, shared initialization.
- `cmd/<command>.go`: command-local flags such as HTTP auth/server settings.
- `internal/infrastructure/config/schema.go`: typed config structs.
- `internal/infrastructure/config/bootstrap.go`: Viper setup, env prefix, config file loading, flag binding.
- `internal/infrastructure/config/config.go`: unmarshal into `config.Config`.
- `config/config.yaml` and `config/example.config.yaml`: YAML values.

## Env Prefix Rule

The env var prefix comes from `appName` in `cmd/root.go`:

```go
strings.ToUpper(strings.ReplaceAll(appName, "-", ""))
```

Examples:

- `go-archetype` becomes `GOARCHETYPE_`.
- `clm-be` becomes `CLMBE_`.

Nested config keys use underscores in env vars. Example: `http.port` becomes `<PREFIX>HTTP_PORT`.

## Adding Config

1. Add fields to `internal/infrastructure/config/schema.go` with `mapstructure` tags.
2. Add a Cobra flag in `cmd/root.go` if shared, or in the owning subcommand if command-specific.
3. Ensure flag names map correctly to config keys: dashes become dots during binding.
4. Add YAML examples to `config/example.config.yaml`; update `config/config.yaml` only when appropriate.
5. Use `*config.Config` after load; avoid direct Viper reads outside config/bootstrap and command setup.
6. Wire the config into adapters or services from `cmd`.

## Pitfalls

- Check whether a flag is bound from `cmd.Flags()` or persistent flags before assuming Viper can read it.
- Keep `mapstructure` names consistent with flag-derived keys.
- Do not introduce package-level config reads inside domain or application code.

