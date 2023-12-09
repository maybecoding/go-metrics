package handlers

import "net/http"

func (c *Handler) ping(w http.ResponseWriter, r *http.Request) {
	err := c.metric.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
