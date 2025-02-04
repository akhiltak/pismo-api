package config

import (
	"log"
	"sync"

	v11env "github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Debug              bool   `env:"DEBUG" envDefault:"false"`
	PostgresDNS        string `env:"POSTGRES_DNS,required,notEmpty"`
	HTTPListenHostPort string `env:"HTTP_LISTEN_HOST_PORT" envDefault:"0.0.0.0:2090"`
}

var instance Config
var instanceOnce sync.Once

// Get returns the config instance
func Get() *Config {
	instanceOnce.Do(func() {
		if err := v11env.Parse(&instance); err != nil {
			log.Fatal("error parsing config")
		}
	})
	return &instance
}
