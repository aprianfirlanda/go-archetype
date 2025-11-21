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
