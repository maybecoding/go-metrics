package config

import (
	"flag"
	"log"
)

type (
	Config struct {
		App
		Sender
	}

	Sender struct {
		Address     string
		Method      string
		Template    string
		IntervalSec int
	}
	App struct {
		CollectIntervalSec int
		SendIntervalSec    int
	}
)

func New() *Config {
	serverAddress := flag.String("a", "localhost:8080", "HTTP server endpoint")
	reportInterval := flag.Int("r", 10, "metric report interval")
	poolInterval := flag.Int("p", 1, "metric pool interval")
	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("undeclared flags provided")
	}
	return &Config{
		App: App{
			CollectIntervalSec: *poolInterval,
			SendIntervalSec:    *reportInterval,
		},

		Sender: Sender{
			Address:  *serverAddress,
			Method:   "POST",
			Template: "http://%s/update/%s/%s/%s",
		},
	}
}
