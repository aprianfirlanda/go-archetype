package config

type Config struct {
	AppName  string
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
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
