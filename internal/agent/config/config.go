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

type (
	Config struct {
		Log     Log
		Sender  Sender `json:"sender"`
		App     App
		CfgFile CfgFile
	}
	App struct {
		CollectInterval time.Duration `default:"2s" flg:"p" flgU:"metric poll interval" env:"POLL_INTERVAL"`
		SendInterval    time.Duration `default:"10s" flg:"r" flgU:"metric report interval" env:"REPORT_INTERVAL"`
	}

	Sender struct {
		Server           string          `default:"localhost:8080" flg:"a" flgU:"HTTP server endpoint" env:"ADDRESS"`
		GRPCServer       string          `flg:"grpc" flgU:"if set instead of http client using gRPC client for send metrics" env:"GRPC"`
		HashKey          string          `flg:"k" flgU:"hash key" env:"KEY"`
		EndpointTemplate string          `default:"%s://%s/update/"`
		CryptoKey        string          `flg:"crypto-key" flgU:"path to certificate" env:"CRYPTO_KEY"`
		RetryIntervals   []time.Duration `default:"1s,3s,5s"`
		NumWorkers       int             `default:"1" flg:"l" flgU:"num workers for send metrics" env:"RATE_LIMIT"`
		IPAddrHeader     string          `default:"X-Real-IP"`
	}

	Log struct {
		Level string `default:"debug" flg:"z" flgU:"log level eg.: debug, error, fatal" env:"LOG_LEVEl"`
	}
	CfgFile struct {
		Path string `flg:"c,config" flgU:"path to config file"  env:"CONFIG"`
	}
)

type FileConfig struct {
	Address        string `json:"address"`
	ReportInterval string `json:"report_interval"`
	PoolInterval   string `json:"pool_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func New() (*Config, error) {
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
		if cfg.Sender.Server == "" {
			cfg.Sender.Server = fCfg.Address
		}
		if cfg.App.SendInterval == 0 {
			rep, err := time.ParseDuration(fCfg.ReportInterval)
			if err != nil {
				logger.Fatal().Err(err).Msg("can't parse report interval")
			}
			cfg.App.SendInterval = rep
		}
		if cfg.App.SendInterval == 0 {
			pool, err := time.ParseDuration(fCfg.PoolInterval)
			if err != nil {
				logger.Fatal().Err(err).Msg("can't parse pool interval")
			}
			cfg.App.SendInterval = pool
		}
		if cfg.Sender.CryptoKey == "" {
			cfg.Sender.CryptoKey = fCfg.CryptoKey
		}
	}
	return cfg, nil
}

func (cfg *Config) LogDebug() {
	logger.Debug().Interface("cfg", cfg).Interface("args", os.Args).Msg("server configuration")
}
func (cfg *Config) UseGRPC() bool {
	return cfg.Sender.GRPCServer != ""
}
