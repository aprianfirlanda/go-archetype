package config

type Config struct {
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
	Services Services `mapstructure:"services"`
}

type Http struct {
	Port int `mapstructure:"port"`
}

type Log struct {
	Format string `mapstructure:"format"`
	Level  string `mapstructure:"level"`
}

type Services struct {
	General Service `mapstructure:"general"`
}

type Service struct {
	BaseURL string `mapstructure:"baseurl"`
	APIKey  string `mapstructure:"apikey"`
}
