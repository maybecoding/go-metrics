package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Handler) metricGetAll(w http.ResponseWriter, r *http.Request) {
	mtr := c.app.GetMetricsAll()

	w.Header().Set("Content-Type", "text/html")

	var result strings.Builder
	result.WriteString(`<h1>Metrics</h1><table><tr><td>#</td><td>Type</td><td>Name</td><td>Value</td></tr><tbody>`)

	for i, metric := range mtr {
		result.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", i, metric.Type, metric.Name, metric.Value))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, result.String())

}
