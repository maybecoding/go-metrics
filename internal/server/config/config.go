// Package config for configuration structs
package config

import (
	"encoding/json"
	"fmt"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/maybecoding/go-metrics.git/pkg/vscfg"
	"os"
	"reflect"
	"time"
)

// Config - root struct
type (
	Config struct {
		Server        Server
		Log           Log
		BackupStorage BackupStorage
		Database      Database
		CfgFile       CfgFile
	}
	// Server - struct for server config
	Server struct {
		Address       string `default:"localhost:8080" flg:"a" flgU:"Endpoint HTTP-server address" env:"ADDRESS"`
		GRPCAddress   string `default:"localhost:9090" flag:"grpc" flgU:"Endpoint gRPC-server address" env:"GRPC_ADDRESS"`
		PprofAddress  string `default:"localhost:8090"`
		HashKey       string `flg:"k" flgU:"hash key" env:"KEY"`
		CryptoKey     string `flg:"crypto-key" flgU:"path to certificate" env:"CRYPTO_KEY"`
		TrustedSubnet string `flg:"t" flgU:"trusted subnet" env:"TRUSTED_SUBNET"`
		IPAddrHeader  string `default:"X-Real-IP"`
	}
	// Log - struct for log config
	Log struct {
		Level string `default:"debug" flg:"l" flgU:"lg level eg.: debug, error, fatal" env:"LOG_LEVEl"`
	}
	// BackupStorage - struct for backup functionality config
	BackupStorage struct {
		Path          string        `default:"/tmp/metric-db.json" flg:"f" flgU:"Storage file path" env:"FILE_STORAGE_PATH"`
		Interval      time.Duration `default:"300s" flg:"i" flgU:"metric backup interval in sec. Default 300 sec" env:"STORE_INTERVAL"`
		IsRestoreOnUp bool          `default:"true" flg:"r" flgU:"Restore data on up" env:"RESTORE"`
	}
	// Database - struct for db config
	Database struct { // postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable
		ConnStr        string          `flg:"d" flgU:"postgres database connection string, if empty - using memory" env:"DATABASE_DSN"`
		RetryIntervals []time.Duration `default:"1s,3s,5s"`
		RunMigrations  bool            `default:"true"`
	}
	CfgFile struct {
		Path string `flg:"c,config" flgU:"path to config file"  env:"CONFIG"`
	}
)

type FileConfig struct {
	Address       string `json:"address"`
	StoreInterval string `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	DatabaseDSN   string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
	Restore       bool   `json:"restore"`
	TrustedSubnet string `json:"trusted_subnet"`
}

// NewConfig - constructor for config structures, reads params from flags and env, env overrides flags
func NewConfig() (*Config, error) {
	cfg := new(Config)
	rCfg := reflect.ValueOf(cfg).Elem()
	// Заполняем значениями по умолчанию, флагами и env
	var fns []vscfg.Fn
	fns = append(fns, vscfg.Tag("default"))
	fns = append(fns, vscfg.Flag("flg", "flgU")...)
	fns = append(fns, vscfg.Env("env"))
	err := vscfg.FillByTags(rCfg, fns...)

	if err != nil {
		return nil, fmt.Errorf("config - New - vscfg.FillByTags: %w", err)
	}
	// Есть указан config-файл пытаемся получить config из него
	// Поскольку этот файл сбоку-припеку и не вписывается в общую модель получения конфигурации заполняем значениями
	// только если они не проставлены ранее.
	// Есть мысли как и это в будущем включить в модель, но структура должна совпадать со структурой используемой конфигурации
	if cfg.CfgFile.Path != "" {
		fCfgB, err := os.ReadFile(cfg.CfgFile.Path)
		if err != nil {
			logger.Fatal().Err(err).Msg("can't read config file")
		}
		var fCfg FileConfig
		err = json.Unmarshal(fCfgB, &fCfg)
		if err != nil {
			logger.Fatal().Err(err).Msg("can't parse config file")
		}

		if cfg.Server.Address == "" {
			cfg.Server.Address = fCfg.Address
		}
		if !cfg.BackupStorage.IsRestoreOnUp && fCfg.Restore {
			cfg.BackupStorage.IsRestoreOnUp = fCfg.Restore
		}
		if cfg.BackupStorage.Interval == 0 {
			storeDur, err := time.ParseDuration(fCfg.StoreInterval)
			if err != nil {
				return nil, err
			}
			cfg.BackupStorage.Interval = storeDur
		}
		if cfg.BackupStorage.Path == "" {
			cfg.BackupStorage.Path = fCfg.StoreFile
		}
		if cfg.Database.ConnStr == "" {
			cfg.Database.ConnStr = fCfg.DatabaseDSN
		}
		if cfg.Server.CryptoKey == "" {
			cfg.Server.CryptoKey = fCfg.CryptoKey
		}
		if cfg.Server.TrustedSubnet == "" {
			cfg.Server.TrustedSubnet = fCfg.TrustedSubnet
		}
	}
	return cfg, nil
}

func (d Database) Use() bool {
	return d.ConnStr != ""
}

func (cfg *Config) LogDebug() {
	logger.Debug().Interface("cfg", cfg).Msg("server configuration")
}
