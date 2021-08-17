package db

import (
	"log"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
    Path      string `required:"true"`
	Username  string
	Pwd		string
}

// NewConfig creates a Config with values from environment variables.
func NewDBConfig(prefix string) *Config {
	c := &Config{}
	if err := envconfig.Process(prefix, c); err != nil {
		log.Fatal(err)
	}
	return c
}