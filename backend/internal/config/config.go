package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port           int      `env:"PORT" envDefault:"8080"`
	DB             string   `env:"DATABASE_URL,required"`
	LogFormat      string   `env:"LOG_FORMAT" envDefault:"text"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"http://localhost:5173"`
	JWTSecret      string   `env:"JWT_SECRET,required"`
}

func Load() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}