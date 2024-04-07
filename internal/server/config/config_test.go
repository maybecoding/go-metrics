package config

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
	"time"
)

func getDefaultConfig() *Config {
	return &Config{
		Server: Server{Address: "localhost:8080", PprofAddress: "localhost:8090", HashKey: "", CryptoKey: "", IPAddrHeader: "X-Real-IP"},
		Log:    Log{Level: "debug"},
		BackupStorage: BackupStorage{
			Interval:      time.Second * 300,
			Path:          "/tmp/metric-db.json",
			IsRestoreOnUp: true,
		},
		Database: Database{
			ConnStr:        "",
			RetryIntervals: []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
			RunMigrations:  true,
		},
	}
}

func TestConfig(t *testing.T) {
	logger.Init("debug")

	tests := []struct {
		name     string
		env      map[string]string
		cfgFile  string
		mutation func(config *Config)
		args     []string
	}{
		{name: "#2 some flags",
			args: []string{"main", "-a", ":8080", "-d", "wrong database"},
			env: map[string]string{"DATABASE_DSN": "database", "LOG_LEVEl": "fatal", "STORE_INTERVAL": "100m",
				"RESTORE": "true", "KEY": "key", "FILE_STORAGE_PATH": "/tmp/tmp.tmp"},
			cfgFile: `{
    "address": "localhost:8080",
    "restore": true,
    "store_interval": "1s",
    "store_file": "/path/to/file.db",
    "database_dsn": "",
    "crypto_key": "/path/to/key.pem",
	"trusted_subnet": "192.168.88.0/24"
}`,
			mutation: func(cfg *Config) {
				cfg.Server.Address = ":8080"
				cfg.Database.ConnStr = "database"
				cfg.Server.CryptoKey = "/path/to/key.pem"
				cfg.Log.Level = "fatal"
				cfg.BackupStorage.Interval = 100 * time.Minute
				cfg.BackupStorage.IsRestoreOnUp = true
				cfg.Server.HashKey = "key"
				cfg.BackupStorage.Path = "/tmp/tmp.tmp"
				cfg.Server.TrustedSubnet = "192.168.88.0/24"
			}},
	}

	for _, tt := range tests {
		println("before check")
		t.Run(tt.name, func(t *testing.T) {
			check(t, tt.args, tt.env, tt.mutation, tt.cfgFile)
		})

	}
}
func check(t *testing.T, args []string, env map[string]string, mutation func(config *Config), cfgFile string) {
	// Задаем файл с логами
	var cfgPath string
	if cfgFile != "" {
		dir, err := os.MkdirTemp("", "tmp")
		require.NoError(t, err)
		cfgPath = path.Join(dir, "cfg.json")
		err = os.WriteFile(cfgPath, []byte(cfgFile), 0666)
		require.NoError(t, err)

		env["CONFIG"] = cfgPath
		defer func() {
			err = os.RemoveAll(dir)
			require.NoError(t, err)
		}()
	}

	// Задаем аргументы
	oldArgs := os.Args
	defer func() {
		println("defer args")
		os.Args = oldArgs
	}()
	os.Args = args

	// Задаем переменные окружения
	oldEnv := make(map[string]string)
	for k, v := range env {
		oldEnv[k] = os.Getenv(k)
		err := os.Setenv(k, v)
		require.NoError(t, err)
	}
	defer func() {
		println("defer os")
		for k, v := range oldEnv {
			err := os.Setenv(k, v)
			require.NoError(t, err)
		}
	}()

	cfg, err := NewConfig()
	require.NoError(t, err)
	cfgExpected := getDefaultConfig()
	cfgExpected.CfgFile.Path = cfgPath
	mutation(cfgExpected)

	require.Equal(t, cfgExpected, cfg)

	if _, ok := env["DATABASE_DSN"]; ok {
		require.Equal(t, true, cfg.Database.Use())
	}

	println("done check")
}
