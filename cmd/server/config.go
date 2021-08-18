package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port             string
	MaxRequestPerSec int
}

// NewConfig creates a Config with values from environment variables.
func NewConfig(prefix string) *Config {
	c := &Config{Port: "8000", MaxRequestPerSec: 1}
	if err := envconfig.Process(prefix, c); err != nil {
		log.Fatal(err)
	}
	return c
}
