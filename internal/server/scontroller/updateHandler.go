package scontroller

import (
	"net/http"
	"strings"
)

func (c *Controller) handleUpdate(w http.ResponseWriter, r *http.Request) {
	// Если метод не POST, значит не обрабатываем запрос
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	url := strings.Split(r.URL.String(), "/")

	// 1/2     /3 /4/5
	// /update/mt/m/v
	if len(url) < 5 || url[1] != "update" {
		http.NotFound(w, r)
		return
	}
	mType := url[2]
	mName := url[3]
	mValue := url[4]

	if err := c.app.UpdateMetric(mType, mName, mValue); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
