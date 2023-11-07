package config

import (
	"flag"
	"log"
	"os"
	"strconv"
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
	// Адрес сервера
	servAddr := flag.String("a", "localhost:8080", "HTTP server endpoint")
	if envServAddr := os.Getenv("ADDRESS"); envServAddr != "" {
		servAddr = &envServAddr
	}

	// Интервал отправки
	repInter := flag.Int("r", 10, "metric report interval")
	if envRepInter := os.Getenv("REPORT_INTERVAL"); envRepInter != "" {
		envRepInterInt, err := strconv.Atoi(envRepInter)
		if err != nil {
			log.Fatal("incorrect REPORT_INTERVAL env value")
		}
		repInter = &envRepInterInt
	}

	poolInter := flag.Int("p", 1, "metric pool interval")
	if envPoolInter := os.Getenv("POLL_INTERVAL"); envPoolInter != "" {
		envPoolInterInt, err := strconv.Atoi(envPoolInter)
		if err != nil {
			log.Fatal("incorrect POLL_INTERVAL env value")
		}
		repInter = &envPoolInterInt
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("undeclared flags provided")
	}
	return &Config{
		App: App{
			CollectIntervalSec: *poolInter,
			SendIntervalSec:    *repInter,
		},

		Sender: Sender{
			Address:  *servAddr,
			Method:   "POST",
			Template: "http://%s/update/%s/%s/%s",
		},
	}
}
