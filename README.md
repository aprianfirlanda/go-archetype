# 🚀 GO Archetype

Initialize the project:
```shell
go mod init go-archetype
```

This project was created using Go version 1.25.4.

Folder Structure
```text
.
├── cmd/
│   ├── root.go          # cobra root cmd, config + logger init
│   └── http.go          # "http" subcommand -> start Fiber server
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── logging/
│   │   └── logging.go
│   ├── database/
│   │   └── gorm.go
│   │
│   ├── domain/          # pure domain & ports (hexagonal core)
│   │   └── customer/
│   │       ├── entity.go          # structs, domain rules
│   │       ├── port_repository.go # interfaces
│   │       └── service.go         # usecase / application service
│   │
│   ├── adapter/
│   │   ├── repository/
│   │   │   └── customer_gorm.go   # implements domain ports using GORM
│   │   └── http/
│   │       └── fiber/
│   │           ├── server.go      # create Fiber app, middlewares, start listen
│   │           ├── router.go      # grouping routes by module
│   │           ├── middleware.go  # CORS, logger, recover, etc.
│   │           └── handler/
│   │               └── customer_handler.go  # Fiber handler -> call usecase
│   │
│   └── pkg/             # optional shared helpers, errors, response wrapper, etc.
│       └── response/
│           └── api_response.go
│
└── go.mod
```

⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻

## 🏗️ Scaffold Application Using Cobra CLI

📥 Install the Cobra library:
```shell
go get -u github.com/spf13/cobra@latest
```

🔧 Install the Cobra CLI generator:
```shell
go install github.com/spf13/cobra-cli@latest
```

📦 Initialize the Cobra project:
```shell
cobra-cli init
```

🧩 Add a command to run the HTTP server:
```shell
cobra-cli add http
```

⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻⸻

## ⚙️ Setup Configuration (Cobra CLI + Viper)

📦 Install Viper

Use the latest version of Viper:
```shell
go get -u github.com/spf13/viper@latest
```

This application uses Cobra CLI and Viper to follow a clean and predictable 12-Factor configuration precedence.
Viper automatically merges values from multiple sources and applies the following order (from highest priority to lowest):

🔝 1. Command-Line Flags

Examples:
```shell
go-archetype http --http-port 3000
go run . http --http-port 3000
go-archetype http --storage-s3-accesskey 3000
go run . http --storage-s3-accesskey 3000
```
Flags must be defined on the root command or the subcommand.

🟩 2. Environment Variables

For this project, environment variables use the prefix derived from the project name:
```shell
GOARCHETYPE_HTTP_PORT=9000
GOARCHETYPE_STORAGE_S3_ACCESSKEY=from_env_var
```

Environment variables override values from the config file.

📄 3. Configuration Files

Viper automatically searches for a configuration file in these locations:
•	Current directory → ./config.yaml
•	Home directory → $HOME/.go-archetype/config.yaml
•	Custom path → via --config internal/config/config.yaml

Example file:
```yaml
http:
  port: 8081
storage:
  s3:
    accessKey: "from_config"
```

🔽 4. Default Values

If flags, environment variables, or config files don’t provide a value, Viper falls back to the default flag value, e.g.:

Subcommand(http) local flag:
```go
httpCmd.Flags().Int("http-port", 8080, "HTTP server port")
```
or

Root persistence flag
```go
rootCmd.PersistentFlags().String("storage-s3-accesskey", "defaultvalue", "access key for S3 storage")
```

This becomes the lowest-priority default.

After that, the value can be called with viper, e.g.:
```go
viper.GetInt("http.port")
viper.GetString("storage.s3.accesskey")
```

But in this app, The config is automatically loaded on Config struct.
Please see the code on persistence pre-run of root command. Update the struct base on config that you add.
This is the struct example:
```go
package config

type Config struct {
	Http Http `mapstructure:"http"`
	Log  Log  `mapstructure:"log"`
}

type Http struct {
	Port int `mapstructure:"port"`
}

type Log struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
}
```

## ✅ Setup Logger using Logrus

This project uses Logrus as the logging library, with full support for:
- Environment variables
- Configuration file (YAML)
- Cobra flags
- Automatic merging of config priority:
  
  flags > env > config file > default

1. Install Logrus
```shell
go get -u github.com/sirupsen/logrus
```

2. Add Logger Configuration Options

Logging can be configured from:

- Environment Variables
```shell
export GOARCHETYPE_LOG_FORMAT=json     # text | json
export GOARCHETYPE_LOG_LEVEL=debug     # trace | debug | info | warn | error | fatal | panic
```

- Configuration File (config.yaml)
```yaml
log:
  format: json   # text or json
  level: debug   # trace | debug | info | warn | error | fatal | panic
```

- Cobra Flags
```shell
go-archetype http --log-format=json --log-level=debug
```

## Setup Web Framework (Fiber)

Install Fiber library
```shell
go get -u github.com/gofiber/fiber/v2
```

The fiber app implements this middleware

- RequestID: The middleware will set context value `requestid` from request header `X-Request-ID`. If does not have that header, it creates its own.
- Recover: if the application that runs under fiber is got panic, the app will not stop.
- CORS: Make the browser block request to api that has different origin.
