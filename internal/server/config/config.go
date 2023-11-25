package config

import (
	"flag"
	"github.com/maybecoding/go-metrics.git/internal/server/logger"
	"os"
)

type (
	Config struct {
		Server Server
		Log    Log
	}

	Server struct {
		Address string
	}

	Log struct {
		Level string
	}
)

func NewConfig() *Config {
	// serverAddress
	serverAddress := flag.String("a", "localhost:8080", "Endpoint HTTP-server address")
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	// logLevel
	logLevel := flag.String("l", "debug", "Log level eg.: debug, error, fatal")
	if envLogLevel := os.Getenv("LOG_LEVEl"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Log.Fatal().Msg("undeclared flags provided")
	}
	return &Config{Server{Address: *serverAddress}, Log{Level: *logLevel}}
}
