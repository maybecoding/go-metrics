package main

import (
	"github.com/maybecoding/go-metrics.git/cmd/server/config"
	sapp "github.com/maybecoding/go-metrics.git/internal/server/sapp"
	"github.com/maybecoding/go-metrics.git/internal/server/scontroller"
	"github.com/maybecoding/go-metrics.git/internal/server/smemstorage"
)

func main() {
	// Получаем конфигурацию приложения
	cfg := config.NewConfig()

	// Создаем хранилище и приложение, приложению даем хранилище
	store := smemstorage.NewMemStorage()
	app := sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	controller := scontroller.New(app, cfg.Server.Address)

	// Запускаем
	controller.Start()

}
