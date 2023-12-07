package config

import (
	"flag"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"os"
	"strconv"
)

type (
	Config struct {
		Server        Server
		Log           Log
		BackupStorage BackupStorage
	}

	Server struct {
		Address string
	}

	Log struct {
		Level string
	}

	BackupStorage struct {
		Interval      int64
		Path          string
		IsRestoreOnUp bool
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
	// storeInterval
	storeIntervalSec := flag.Int64("i", 300, "Metric backup interval in sec. Default 300 sec")
	if envStoreIntervalSec := os.Getenv("STORE_INTERVAL"); envStoreIntervalSec != "" {
		parsed, err := strconv.ParseInt(envStoreIntervalSec, 10, 64)
		if err != nil {
			logger.Log.Panic().Err(err).Msg("can't parse int from STORE_INTERVAL in env")
		}
		storeIntervalSec = &parsed
	}

	// fileStoragePath
	fileStoragePath := flag.String("f", "/tmp/metrics-db.json", "Storage file path")
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

	flag.Parse()
	if len(flag.Args()) > 0 {
		logger.Log.Fatal().Msg("undeclared flags provided")
	}
	return &Config{
		Server: Server{Address: *serverAddress},
		Log:    Log{Level: *logLevel},
		BackupStorage: BackupStorage{
			Interval:      *storeIntervalSec,
			Path:          *fileStoragePath,
			IsRestoreOnUp: *isRestoreOnUp,
		},
	}
}
