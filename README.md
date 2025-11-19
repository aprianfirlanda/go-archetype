# 🚀 GO Archetype

Initialize the project:
```shell
go mod init go-archetype
```

This project was created using Go version 1.25.4.

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
go-archetype http --port 3000
go run main.go http --port 3000
```
Flags must be defined on the root command or the subcommand.

🟩 2. Environment Variables

For this project, environment variables use the prefix derived from the project name:
```shell
GOARCHETYPE_PORT=9000
```

Environment variables override values from the config file.

📄 3. Configuration Files

Viper automatically searches for a configuration file in these locations:
•	Current directory → ./config.yaml
•	Home directory → $HOME/.go-archetype/config.yaml
•	Custom path → via --config config/config.yaml

Example file:
```yaml
port: 8081
```

🔽 4. Default Values

If flags, environment variables, or config files don’t provide a value, Viper falls back to the default flag value, e.g.:
```go
httpCmd.Flags().Int("port", 8080, "HTTP server port")
```
This becomes the lowest-priority default.



