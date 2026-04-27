# 🚀 GO Archetype

Initialize the project:
```shell
go mod init go-archetype
```

This project was created using Go version 1.25.4.

Folder Structure
```text
go-archetype/
├── bin/                          # Compiled binaries (optional, local use)
│
├── cmd/                          # CLI entry points (Cobra commands)
│   ├── root.go                   # Root command & global flags
│   ├── http.go                   # Start HTTP server command
│   └── migrate.go                # Run database migration command
│
├── config/                       # Application configuration files
│   ├── config.yaml               # Main configuration
│   └── example.config.yaml       # Example / template configuration
│
├── internal/                     # Private application code (not importable outside)
│
│   ├── adapters/                 # Inbound adapters (how the app is called)
│   │   └── http/                 # HTTP adapter (Fiber)
│   │       ├── context/          # HTTP context helpers (logger, request-scoped data)
│   │       │   └── logger.go
│   │       │
│   │       ├── docs/             # API documentation (Swagger / OpenAPI)
│   │       │   ├── docs.go
│   │       │   ├── swagger.json
│   │       │   └── swagger.yaml
│   │       │
│   │       ├── dto/              # Data Transfer Objects for HTTP
│   │       │   ├── request/      # Incoming HTTP request payloads
│   │       │   │   └── task/
│   │       │   └── response/     # Outgoing HTTP response payloads
│   │       │       ├── task/
│   │       │       ├── common.go
│   │       │       ├── error.go
│   │       │       ├── helper.go
│   │       │       ├── meta.go
│   │       │       └── success.go
│   │       │
│   │       ├── handler/          # HTTP handlers (request → usecase → response)
│   │       │   ├── demo/         # Demo / example endpoints
│   │       │   └── task/         # Task-related HTTP handlers
│   │       │
│   │       ├── middleware/       # HTTP middleware (auth, logging, recovery, etc.)
│   │       │
│   │       ├── router/           # HTTP route definitions
│   │       │   └── router.go
│   │       │
│   │       ├── server/           # HTTP server setup (Fiber instance)
│   │       │   └── fiber.go
│   │       │
│   │       └── validation/       # Request validation logic
│   │           ├── field.go
│   │           ├── message.go
│   │           └── validator.go
│
│   ├── application/              # Application layer (use cases)
│   │   ├── health/               # Health check use cases
│   │   │   └── service/
│   │   │
│   │   └── task/                 # Task-related use cases
│   │       ├── command/          # Write operations (create, update, delete)
│   │       ├── query/            # Read operations (list, filter)
│   │       ├── result/           # Usecase result models
│   │       └── service/          # Orchestration logic
│
│   ├── domain/                   # Core business rules (pure Go)
│   │   ├── auth/                 # Authentication-related domain logic
│   │   ├── identity/             # Identity concepts (public ID, etc.)
│   │   └── task/                 # Task domain (entities, status, errors)
│
│   ├── infrastructure/           # Outbound adapters & framework implementations
│   │   ├── config/               # Configuration loading & schema
│   │   │
│   │   ├── db/                   # Database connection & pooling
│   │   │
│   │   ├── logging/              # Logging implementation (Logrus, fields, components)
│   │   │
│   │   └── persistance/          # Data persistence implementations
│   │       └── gorm/             # GORM-based repository implementations
│   │           ├── migrate/      # Database migration tooling (Goose)
│   │           ├── task/         # Task repository implementation
│   │           ├── uow.go        # Unit of Work implementation
│   │           ├── uow_tx.go     # Transaction handling
│   │           ├── bootstrap.go  # Persistence wiring
│   │           └── pinger.go     # Database health check
│
│   ├── ports/                    # Hexagonal ports (interfaces / contracts)
│   │   ├── input/                # Inbound ports (usecase interfaces)
│   │   │   ├── task_service.go
│   │   │   └── health_service.go
│   │   └── output/               # Outbound ports (DB, UoW, external systems)
│   │       ├── repository.go
│   │       ├── uow.go
│   │       └── db_pinger.go
│
│   ├── bootstrap/                # Dependency wiring (compose adapters & usecases)
│   │   └── http_app.go
│
│   └── pkg/                      # Shared utilities (non-domain helpers)
│       └── apperror/
│           ├── error.go
│           └── helpers.go
│
├── migrations/                   # SQL migration files
│   └── 20260104053845_create_tasks_table.sql
│
├── test/                         # Test resources
│   └── http/
│       └── request/             # HTTP request collections (manual testing)
│
├── compose.yaml                  # Docker Compose for local development
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
├── go.sum                        # Go dependencies checksum
├── README.md                     # Project documentation
└── LICENSE                       # License
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

The Fiber application in this project uses several global middlewares to ensure observability, stability, and browser compatibility.

🔹 Middleware Implemented

1. RequestID

   Adds a unique request ID to every incoming request.
   - Reads X-Request-ID from the request header if provided.
   - If the header is missing, it generates a new UUID.
   - Stores the value in c.Locals("requestid") so it can be used by:
   - Logging middleware
   - Error recover middleware
   - Handlers

   This improves tracing and debugging across distributed systems.
                      
2. Logging (custom middleware)

   Logs every request using the application’s shared Logrus logger. The log entry includes structured fields such as:
   - request_id
   - status
   - method
   - path
   - latency
   - ip

   This makes the log output consistent and easy to search in log aggregation systems (e.g., Loki, Elasticsearch, OpenSearch).

3. Recover (custom Fiber wrapper)

   Prevents the application from crashing due to panics.
   - Catches any panic in handlers or other middleware
   - Logs the error and full stack trace using Logrus
   - Returns a safe 500 Internal Server Error response
   - Ensures the server continues running even when unexpected errors occur

   This provides fault tolerance and better observability during debugging.

4. CORS

   Enables Cross-Origin Resource Sharing.
   - Allows frontend applications on different domains/ports to access the API
   - Prevents browser CORS errors during local development
   - Configurable for production environments

   This is required when the frontend runs on localhost:5173 (Vite) and the backend runs on another port.

5. Health Check (custom middleware) (if included)

   Provides lightweight liveness/readiness endpoints for:
   - Kubernetes
   - Docker health checks
   - Load balancers

   Example endpoints:
   - GET /live → checks if the process is alive
   - GET /ready → checks if the service is ready (e.g., DB connection)

6. Metrics (Fiber Monitor)

   The application exposes a built-in metrics dashboard using Fiber’s monitor middleware.
   - Accessible via: 
   
     `GET /metrics`
   
     Using Browser or curl `-H "Accept: application/json"`
   - Displays:
     - Total requests
     - Status code counts
     - Average latency
     - Memory usage
     - Goroutine count
     - Uptime
     - Per-route performance

   This endpoint is useful for basic diagnostics and local performance monitoring.

7. Auth API Key

   Create random text for API Key:
   ```shell
   openssl rand -base64 32
   ```

   If required, the API Key should send on request header `Authorization: Bearer <API_KEY>`.
   To Generate API Key, use this command.

add validation library
```shell
go get -u github.com/go-playground/validator/v10 
```

add swagger library for fiber
```shell
go get -u github.com/swaggo/fiber-swagger
go get -u github.com/swaggo/swag
```

generate swagger docs
```shell
swag init \
  -g cmd/http.go \
  -o internal/adapters/http/docs
```

access swagger ui at
```http://localhost:{port}/swagger/index.html
```

## Setup Database (GORM)

install dependency
```shell
go get -u gorm.io/gorm@lates
go get -u gorm.io/driver/postgres@latest   # or mysql / sqlite / sqlserver
```

setup config using the cobra flag
```shell
go run . http \
--db-host=localhost \
--db-port=5432 \
--db-user=app \
--db-password=change_me \
--db-name=app \
--db-sslmode=disable \
--db-timezone=UTC \
--db-maxopenconns=25 \
--db-maxidleconns=25 \
--db-connmaxlifetime=10h \
--db-connmaxidletime=15m \
--db-loglevel=warn
```

setup config using env
```shell
export GOARCHETYPE_DB_HOST=localhost
export GOARCHETYPE_DB_PORT=5432
export GOARCHETYPE_DB_USER=app
export GOARCHETYPE_DB_PASSWORD=change_me
export GOARCHETYPE_DB_NAME=app
export GOARCHETYPE_DB_SSLMODE=disable
export GOARCHETYPE_DB_TIMEZONE=UTC
export GOARCHETYPE_DB_MAXOPENCONNS=25
export GOARCHETYPE_DB_MAXIDLECONNS=25
export GOARCHETYPE_DB_CONNMAXLIFETIME=10h
export GOARCHETYPE_DB_CONNMAXIDLETIME=15m
export GOARCHETYPE_DB_LOGLEVEL=warn
```

setup config using config file
```yaml
db:
  host: localhost
  port: 5432
  user: app
  password: change_me
  name: app
  sslMode: disable
  timezone: UTC
  maxOpenConns: 25
  maxIdleConns: 25
  connMaxLifetime: 10h
  connMaxIdleTime: 15m
  logLevel: warn
```

add lib for generate uuid v7
```shell
go get -u github.com/google/uuid
```

## Setup Migration (GOOSE - GORM)
install goose cli
```shell
go get -u github.com/pressly/goose/v3
```

This project uses Goose for database migrations, integrated directly into the application CLI (Cobra).

All migration commands are executed via:
```shell
go run . migrate <command>
```
⚠️ Important
All migrate commands require a valid database connection, since the database is initialized when the application starts.

⸻

📁 Migration Directory

Migration files are stored in:
```text
migrations/
```
Example:
```text
migrations/
├── 20260104120000_init_tasks.sql
├── 20260104123000_add_index.sql
```

✨ Migration Commands

1️⃣ Create a New Migration
```shell
go run . migrate create <name>
```
Example:
```shell
go run . migrate create init_tasks
```
By default, this creates a SQL migration.

Migration Type Options
```shell
go run . migrate create add_users_table --type sql
go run . migrate create add_indexes --type go
```

2️⃣ Apply All Pending Migrations

Runs all migrations that have not yet been applied:
```shell
go run . migrate up
```

3️⃣ Migrate Up to a Specific Version
```shell
go run . migrate up-to <version>
```
Example:
```shell
go run . migrate up-to 20260104120000
```

4️⃣ Roll Back the Last Migration
```shell
go run . migrate down
```

5️⃣ Roll Back to a Specific Version
```shell
go run . migrate down-to <version>
```
Example:
```shell
go run . migrate down-to 20260104120000
```

6️⃣ Show Migration Status

Displays all migrations and their current status (applied / pending):
```shell
go run . migrate status
```

7️⃣ Show Current Database Version
```shell
go run . migrate version
```
Example output:
```shell
20260104120000
```


TODO:
- [ ] implement http context get user's info
- [ ] use domain auth custom claims for user's info if possible
- [ ] create example extract username, permission, validate permission from c.Locals("user"). Combine with custom claims if posible

