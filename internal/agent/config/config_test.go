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
		App: App{
			CollectInterval: 2 * time.Second,
			SendInterval:    10 * time.Second,
		},

		Sender: Sender{
			Server:           "localhost:8080",
			EndpointTemplate: "%s://%s/update/",
			RetryIntervals:   []time.Duration{time.Second, 3 * time.Second, 5 * time.Second},
			HashKey:          "",
			NumWorkers:       1,
			CryptoKey:        "",
		},

		Log: Log{
			Level: "debug",
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
			args: []string{"main", "-a", ":8080", "-z", "wrong"},
			env: map[string]string{"LOG_LEVEl": "fatal", "REPORT_INTERVAL": "1000m", "POLL_INTERVAL": "3h",
				"KEY": "key", "RATE_LIMIT": "100",
			},
			cfgFile: `{
    "address": "localhost:8080",
    "restore": true,
    "store_interval": "1s",
    "store_file": "/path/to/file.db",
    "database_dsn": "",
    "crypto_key": "/path/to/key.pem"
} `,
			mutation: func(cfg *Config) {
				cfg.Sender.Server = ":8080"
				cfg.Log.Level = "fatal"
				cfg.Sender.CryptoKey = "/path/to/key.pem"
				cfg.App.SendInterval = 1000 * time.Minute
				cfg.App.CollectInterval = 3 * time.Hour
				cfg.Sender.HashKey = "key"
				cfg.Sender.NumWorkers = 100
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
	if cfgFile != "" {
		dir, err := os.MkdirTemp("", "tmp")
		require.NoError(t, err)
		cfgPath := path.Join(dir, "cfg.json")
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

	cfg := New()
	cfgExpected := getDefaultConfig()
	mutation(cfgExpected)

	require.Equal(t, cfgExpected, cfg)
	println("done check")
}
