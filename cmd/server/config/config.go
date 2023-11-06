package config

import (
	"flag"
	"log"
	"os"
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
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("undeclared flags provided")
	}
	return &Config{Server{Address: *serverAddress}}
}
