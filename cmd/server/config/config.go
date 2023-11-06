package config

import (
	"flag"
	"log"
)

type (
	Config struct {
		Server
	}

	Server struct {
		Address string
	}
)

func NewConfig() *Config {
	serverAddress := flag.String("a", "localhost:8080", "Endpoint HTTP-server address")
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("undeclared flags provided")
	}
	return &Config{Server{Address: *serverAddress}}
}
