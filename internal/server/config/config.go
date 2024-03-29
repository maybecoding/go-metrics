// Package config for configuration structs
package config

import (
	"flag"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"strconv"
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
	}
	// Log - struct for log config
	Log struct {
		Level string
	}
	// BackupStorage - struct for backup functionality config
	BackupStorage struct {
		Path          string
		Interval      int64
		IsRestoreOnUp bool
	}
	// Database - struct for db config
	Database struct {
		ConnStr        string
		RetryIntervals []time.Duration
		RunMigrations  bool
	}
)

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
	storeIntervalSec := flag.Int64("i", 300, "metric backup interval in sec. Default 300 sec")
	if envStoreIntervalSec := os.Getenv("STORE_INTERVAL"); envStoreIntervalSec != "" {
		parsed, err := strconv.ParseInt(envStoreIntervalSec, 10, 64)
		if err != nil {
			logger.Fatal().Err(err).Msg("can't parse int from STORE_INTERVAL in env")
		}
		storeIntervalSec = &parsed
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

	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Fatal().Msg("undeclared flags provided")
	}
	cfg := &Config{
		Server: Server{Address: *serverAddress, PprofAddress: "localhost:8090", HashKey: *key},
		Log:    Log{Level: *logLevel},
		BackupStorage: BackupStorage{
			Interval:      *storeIntervalSec,
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
