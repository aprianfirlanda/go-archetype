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
	SSLMode  string `mapstructure:"sslmode"`
	TimeZone string `mapstructure:"timezone"`

	MaxOpenConns    int           `mapstructure:"maxopenconns"`
	MaxIdleConns    int           `mapstructure:"maxidleconns"`
	ConnMaxLifetime time.Duration `mapstructure:"connmaxlifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"connmaxidletime"`

	LogLevel string `mapstructure:"loglevel"`
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
