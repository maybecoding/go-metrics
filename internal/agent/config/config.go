package config

import (
	"flag"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"log"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		App    App
		Sender Sender
		Log    Log
	}

	Sender struct {
		Address           string
		Method            string
		Template          string
		JSONEndpoint      string
		JSONBatchEndpoint string
		IntervalSec       int
		RetryIntervals    []time.Duration
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
	pollInter := flag.Int("p", 2, "metric poll interval")
	if envPollInter := os.Getenv("POLL_INTERVAL"); envPollInter != "" {
		envPollInterInt, err := strconv.Atoi(envPollInter)
		if err != nil {
			log.Fatal("incorrect POLL_INTERVAL env value")
		}
		repInter = &envPollInterInt
	}

	// Уровень логирования
	logLevel := flag.String("l", "debug", "lg level eg.: debug, error, fatal")
	if envLogLevel := os.Getenv("LOG_LEVEl"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Fatal().Msg("undeclared flags provided")
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		log.Fatal("undeclared flags provided")
	}
	return &Config{
		App: App{
			CollectIntervalSec: *pollInter,
			SendIntervalSec:    *repInter,
		},

		Sender: Sender{
			Address:           *servAddr,
			Method:            "POST",
			Template:          "http://%s/update/%s/%s/%s",
			JSONEndpoint:      "http://%s/update/",
			JSONBatchEndpoint: "http://%s/updates/",
			RetryIntervals:    []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
		},

		Log: Log{
			Level: *logLevel,
		},
	}
}
