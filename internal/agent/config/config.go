package config

import (
	"encoding/json"
	"flag"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"strconv"
	"time"
)

type (
	Config struct {
		Log    Log
		Sender Sender `json:"sender"`
		App    App
	}
	App struct {
		CollectInterval time.Duration
		SendInterval    time.Duration
	}

	Sender struct {
		Server           string `json:"address"`
		Method           string
		HashKey          string
		EndpointTemplate string
		CryptoKey        string
		RetryIntervals   []time.Duration
		NumWorkers       int
	}

	Log struct {
		Level string
	}
)

type FileConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PoolInterval   string `json:"pool_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func New() *Config {
	// Адрес сервера
	servAddr := flag.String("a", "localhost:8080", "HTTP server endpoint")
	if envServAddr := os.Getenv("ADDRESS"); envServAddr != "" {
		servAddr = &envServAddr
	}

	// Интервал отправки
	repInter := flag.String("r", "10s", "metric report interval")
	if envRepInter := os.Getenv("REPORT_INTERVAL"); envRepInter != "" {
		repInter = &envRepInter
	}

	// Интервал сборки
	pollInter := flag.String("p", "2s", "metric poll interval")
	if envPollInter := os.Getenv("POLL_INTERVAL"); envPollInter != "" {
		pollInter = &envPollInter
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
		if err == nil {
			numWorkers = &num
		}
	}

	// Публичный ключ
	cryptoKey := flag.String("crypto-key", "", "path to certificate")
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cryptoKey = &envCryptoKey
	}

	// Файл конфигурации
	var configFile string
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.StringVar(&configFile, "config", "", "path to config file")
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Fatal().Msg("undeclared flags provided")
	}
	// Есть указан config-файл пытаемся получить config из него
	if configFile != "" {
		fCfgB, err := os.ReadFile(configFile)
		if err != nil {
			logger.Fatal().Err(err).Msg("can't read config file")
		}
		var fCfg FileConfig
		err = json.Unmarshal(fCfgB, &fCfg)
		if err != nil {
			logger.Fatal().Err(err).Msg("can't parse config file")
		}

		if *servAddr == "" {
			*servAddr = fCfg.Address
		}
		if *repInter == "" {
			*repInter = fCfg.ReportInterval
		}
		if *pollInter == "" {
			*pollInter = fCfg.PoolInterval
		}
		if *cryptoKey == "" {
			*cryptoKey = fCfg.CryptoKey
		}
	}

	pool, err := time.ParseDuration(*pollInter)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't parse pool interval")
	}
	rep, err := time.ParseDuration(*repInter)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't parse report interval")
	}

	cfg := &Config{
		App: App{
			CollectInterval: pool,
			SendInterval:    rep,
		},

		Sender: Sender{
			Server:           *servAddr,
			EndpointTemplate: "%s://%s/update/",
			RetryIntervals:   []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
			HashKey:          *key,
			NumWorkers:       *numWorkers,
			CryptoKey:        *cryptoKey,
		},

		Log: Log{
			Level: *logLevel,
		},
	}
	logger.Debug().Str("key", *key).Msg("agent configuration")
	return cfg
}

func (cfg *Config) LogDebug() {
	logger.Debug().Interface("cfg", cfg).Interface("args", os.Args).Msg("server configuration")
}
