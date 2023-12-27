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
		Sender Sender `json:"sender"`
		Log    Log
	}
	App struct {
		CollectIntervalSec int
		SendIntervalSec    int
	}

	Sender struct {
		Server           string `json:"address"`
		Method           string
		HashKey          string
		EndpointTemplate string
		RetryIntervals   []time.Duration
		NumWorkers       int
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
	logLevel := flag.String("z", "debug", "lg level eg.: debug, error, fatal")
	if envLogLevel := os.Getenv("LOG_LEVEl"); envLogLevel != "" {
		logLevel = &envLogLevel
	}

	// Ключ хеширования
	key := flag.String("k", "", "hash key")
	if envKey := os.Getenv("KEY"); envKey != "" {
		key = &envKey
	}

	// Число одновременных отправок метрик
	numWorkers := flag.Int("l", 1, "num workers for send metrics")
	if envNumWorkers := os.Getenv("RATE_LIMIT"); envNumWorkers != "" {
		num, err := strconv.Atoi(envNumWorkers)
		if err != nil {
			numWorkers = &num
		}
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Fatal().Msg("undeclared flags provided")
	}

	cfg := &Config{
		App: App{
			CollectIntervalSec: *pollInter,
			SendIntervalSec:    *repInter,
		},

		Sender: Sender{
			Server:           *servAddr,
			EndpointTemplate: "http://%s/update/",
			RetryIntervals:   []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
			HashKey:          *key,
			NumWorkers:       *numWorkers,
		},

		Log: Log{
			Level: *logLevel,
		},
	}
	//logger.Debug().Str("key", *key).Msg("agent configuration")
	return cfg
}

func (cfg *Config) LogDebug() {
	logger.Debug().Interface("cfg", cfg).Interface("args", os.Args).Msg("server configuration")
}

func (a App) CollectInterval() time.Duration {
	return time.Duration(a.CollectIntervalSec) * time.Second
}
func (a App) SendInterval() time.Duration {
	return time.Duration(a.SendIntervalSec) * time.Second
}
