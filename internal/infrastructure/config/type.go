package config

import "time"

type Config struct {
	AppName  string
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
	DB       Database `mapstructure:"db"`
	JWT      JWT      `mapstructure:"jwt"`
	Services Services `mapstructure:"services"`
}

type Http struct {
	Port int `mapstructure:"port"`
}

type Log struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`  // e.g. "disable", "require"
	TimeZone string `mapstructure:"timezone"` // e.g. "UTC"

	// Connection pool settings (optional, with sensible defaults if zero).
	MaxOpenConns    int           `mapstructure:"maxopenconns"`    // default: 25
	MaxIdleConns    int           `mapstructure:"maxidleconns"`    // default: 25
	ConnMaxLifetime time.Duration `mapstructure:"connmaxlifetime"` // default: 1h
	ConnMaxIdleTime time.Duration `mapstructure:"connmaxidletime"` // default: 15m

	// GORM logger level (optional).
	LogLevel string `mapstructure:"loglevel"` // Typical values: "silent", "error", "warn", "info". default: nil(warn)
}

type JWT struct {
	Secret string `mapstructure:"secret"`
}

type Services struct {
	General Service `mapstructure:"general"`
}

type Service struct {
	BaseURL string `mapstructure:"baseurl"`
	APIKey  string `mapstructure:"apikey"`
}
