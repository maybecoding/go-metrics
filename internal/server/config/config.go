// Package config for configuration structs
package config

import (
	"encoding/json"
	"flag"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"time"
)

// Config - root struct
type (
	Config struct {
		Server        Server
		Log           Log
		BackupStorage BackupStorage
		Database      Database
	}
	// Server - struct for server config
	Server struct {
		Address      string
		PprofAddress string
		HashKey      string
		CryptoKey    string
	}
	// Log - struct for log config
	Log struct {
		Level string
	}
	// BackupStorage - struct for backup functionality config
	BackupStorage struct {
		Path          string
		Interval      time.Duration
		IsRestoreOnUp bool
	}
	// Database - struct for db config
	Database struct {
		ConnStr        string
		RetryIntervals []time.Duration
		RunMigrations  bool
	}
)

type FileConfig struct {
	Address       string `json:"address"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	Restore       bool   `json:"restore"`
}

// NewConfig - constructor for config structures, reads params from flags and env, env overrides flags
func NewConfig() *Config {
	// serverAddress
	serverAddress := flag.String("a", "localhost:8080", "Endpoint HTTP-server address")
	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		serverAddress = &envServerAddress
	}
	// logLevel
	logLevel := flag.String("l", "debug", "lg level eg.: debug, error, fatal")
	if envLogLevel := os.Getenv("LOG_LEVEl"); envLogLevel != "" {
		logLevel = &envLogLevel
	}
	// storeInterval
	storeInterval := flag.String("i", "300s", "metric backup interval in sec. Default 300 sec")
	if envStoreIntervalSec := os.Getenv("STORE_INTERVAL"); envStoreIntervalSec != "" {
		storeInterval = &envStoreIntervalSec
	}

	// fileStoragePath
	fileStoragePath := flag.String("f", "/tmp/metric-db.json", "Storage file path")
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = &envFileStoragePath
	}

	// isRestoreOnUp
	isRestoreOnUp := flag.Bool("r", true, "Restore data on up")
	if envIsRestoreOnUp := os.Getenv("RESTORE"); envIsRestoreOnUp != "" {
		var envIsRestoreOnUpBool bool
		if envIsRestoreOnUp == "true" {
			envIsRestoreOnUpBool = true
		}
		isRestoreOnUp = &envIsRestoreOnUpBool
	}

	// databaseConnStr
	//databaseConnStr := flag.String("d", "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable", "postgres database connection string, if empty - using")
	databaseConnStr := flag.String("d", "", "postgres database connection string, if empty - using")
	if envDatabaseConnStr := os.Getenv("DATABASE_DSN"); envDatabaseConnStr != "" {
		databaseConnStr = &envDatabaseConnStr
	}

	// Ключ хеширования
	key := flag.String("k", "", "hash key")
	if envKey := os.Getenv("KEY"); envKey != "" {
		key = &envKey
	}

	cryptoKey := flag.String("crypto-key", "", "path to certificate")
	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cryptoKey = &envCryptoKey
	}

	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Fatal().Msg("undeclared flags provided")
	}

	// Файл конфигурации
	var configFile string
	flag.StringVar(&configFile, "c", "", "path to config file")
	flag.StringVar(&configFile, "config", "", "path to config file")
	if envConfigFile := os.Getenv("CONFIG"); envConfigFile != "" {
		configFile = envConfigFile
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

		if *serverAddress == "" {
			*serverAddress = fCfg.Address
		}
		if !*isRestoreOnUp && fCfg.Restore {
			*isRestoreOnUp = fCfg.Restore
		}
		if *storeInterval == "" {
			*storeInterval = fCfg.StoreInterval
		}
		if *fileStoragePath == "" {
			*fileStoragePath = fCfg.StoreFile
		}
		if *databaseConnStr == "" {
			*databaseConnStr = fCfg.DatabaseDSN
		}
		if *cryptoKey == "" {
			*cryptoKey = fCfg.CryptoKey
		}
	}

	storeDur, err := time.ParseDuration(*storeInterval)
	if err != nil {
		logger.Fatal().Err(err).Msg("can't parse pool interval")
	}

	cfg := &Config{
		Server: Server{Address: *serverAddress, PprofAddress: "localhost:8090", HashKey: *key, CryptoKey: *cryptoKey},
		Log:    Log{Level: *logLevel},
		BackupStorage: BackupStorage{
			Interval:      storeDur,
			Path:          *fileStoragePath,
			IsRestoreOnUp: *isRestoreOnUp,
		},
		Database: Database{
			ConnStr:        *databaseConnStr,
			RetryIntervals: []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
			RunMigrations:  true,
		},
	}

	return cfg
}

func (d Database) Use() bool {
	return d.ConnStr != ""
}

func (cfg *Config) LogDebug() {
	logger.Debug().Interface("cfg", cfg).Msg("server configuration")
}
