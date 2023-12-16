package handlers

import "net/http"

func (c *Handler) ping(w http.ResponseWriter, r *http.Request) {
	if c.health.Check() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
