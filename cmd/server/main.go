package main

import (
	sapp "github.com/maybecoding/go-metrics.git/internal/server/app"
	"github.com/maybecoding/go-metrics.git/internal/server/config"
	"github.com/maybecoding/go-metrics.git/internal/server/controller"
	"github.com/maybecoding/go-metrics.git/internal/server/memstorage"
)

func main() {
	// Получаем конфигурацию приложения
	cfg := config.NewConfig()

	// Создаем хранилище и приложение, приложению даем хранилище
	var store sapp.Storage = smemstorage.NewMemStorage()
	app := sapp.New(store)

	// Создаем контроллер и вверяем ему приложение
	contr := controller.New(app, cfg.Server.Address)

	// Запускаем
	contr.Start()

}
