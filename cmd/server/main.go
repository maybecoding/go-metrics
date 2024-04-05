package main

import (
	"fmt"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
)

func main() {
	printInfo()
	// Получаем конфигурацию приложения
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("failed to parse config: %s", err.Error())
		return
	}
	sapp.New(cfg).
		Init().
		InitHandler().
		Run()
}
