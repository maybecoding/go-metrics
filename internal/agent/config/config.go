package config

import (
	"flag"
	"github.com/maybecoding/go-metrics.git/internal/agent/logger"
	"log"
	"os"
	"strconv"
)

type (
	Config struct {
		App    App
		Sender Sender
		Log    Log
	}

	Sender struct {
		Address      string
		Method       string
		Template     string
		JSONEndpoint string
		IntervalSec  int
	}
	App struct {
		CollectIntervalSec int
		SendIntervalSec    int
	}
	Log struct {
		Level string
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

	// Интервал сборки
	poolInter := flag.Int("p", 3, "metric pool interval")
	if envPoolInter := os.Getenv("POLL_INTERVAL"); envPoolInter != "" {
		envPoolInterInt, err := strconv.Atoi(envPoolInter)
		if err != nil {
			log.Fatal("incorrect POLL_INTERVAL env value")
		}
		repInter = &envPoolInterInt
	}

	// Уровень логирования
	logLevel := flag.String("l", "debug", "Log level eg.: debug, error, fatal")
	if envLogLevel := os.Getenv("LOG_LEVEl"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Log.Fatal().Msg("undeclared flags provided")
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
			Address:      *servAddr,
			Method:       "POST",
			Template:     "http://%s/update/%s/%s/%s",
			JSONEndpoint: "http://%s/update",
		},

		Log: Log{
			Level: *logLevel,
		},
	}
}
