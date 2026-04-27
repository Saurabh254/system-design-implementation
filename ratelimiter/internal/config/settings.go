package config

import "fmt"

type Config struct {
	Port int
	Host string
}

func Load() *Config {
	return &Config{
		Port: getEnvInt("PORT", 8080),
		Host: getEnv("HOST", "0.0.0.0"),
	}
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
