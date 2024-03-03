package main

import (
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
)

func main() {
	// Получаем конфигурацию приложения
	cfg := config.NewConfig()
	sapp.New(cfg).
		Init().
		InitHandler().
		Run()
}
