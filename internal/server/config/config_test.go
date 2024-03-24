package config

import (
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func getDefaultConfig() *Config {
	return &Config{
		Server: Server{Address: "localhost:8080", PprofAddress: "localhost:8090", HashKey: "", CryptoKey: ""},
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
		mutation func(config *Config)
		args     []string
	}{
		{name: "#2 some flags",
			args: []string{"main", "-a", ":8080", "-d", "wrong database"},
			env:  map[string]string{"DATABASE_DSN": "database"},
			mutation: func(cfg *Config) {
				cfg.Server.Address = ":8080"
				cfg.Database.ConnStr = "database"
			}},
	}

	for _, tt := range tests {
		println("before check")
		t.Run(tt.name, func(t *testing.T) {
			check(t, tt.args, tt.env, tt.mutation)
		})

	}
}
func check(t *testing.T, args []string, env map[string]string, mutation func(config *Config)) {
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

	cfg := NewConfig()
	cfgExpected := getDefaultConfig()
	mutation(cfgExpected)

	require.Equal(t, cfgExpected, cfg)
	println("done check")
}
