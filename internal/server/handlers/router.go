package handlers

import (
	"crypto/sha256"
	"github.com/go-chi/chi/v5"
	"github.com/maybecoding/go-metrics.git/internal/server/handlers/middleware"
	"github.com/maybecoding/go-metrics.git/pkg/compress"
	"github.com/maybecoding/go-metrics.git/pkg/hashcheck"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
)

func (c *Handler) GetRouter() chi.Router {

	r := chi.NewRouter()

	// Подключаем логгер
	r.Use(logger.Handler)

	// Проверяем что узел от которого принимаем метрики находится в доверенной зоне
	r.Use(middleware.CheckTrustedSubnet(c.trustedSubNet, c.cfg.IPAddrHeader))

	// Подключаем проверку хэшей
	hashCheck := hashcheck.New(sha256.New, c.cfg.HashKey, "HashSHA256")
	//r.Use(hashCheck.Handler)

	// Установка значений
	r.Post("/update/{type}/{name}/{value}", c.metricUpdate)
	r.Post("/update/", compress.HandlerFuncReader(compress.HandlerFuncWriter(c.metricUpdateJSON, compress.BestSpeed)))

	r.Post("/updates/", hashCheck.HandlerFunc(compress.HandlerFuncReader(c.metricUpdateAllJSON)))

	// Получение значений
	r.Get("/value/{type}/{name}", c.metricGet)
	r.Post("/value/", compress.HandlerFuncReader(compress.HandlerFuncWriter(c.metricGetJSON, compress.BestSpeed)))

	// Отчет по метрикам
	r.Get("/", compress.HandlerFuncWriter(c.metricGetAll, compress.BestSpeed))

	// ping
	r.Get("/ping", c.ping)

	return r
}
