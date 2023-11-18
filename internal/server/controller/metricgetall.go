package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Controller) metricGetAll(w http.ResponseWriter, r *http.Request) {
	mtr := c.app.GetMetricsAll()
	var result strings.Builder

	result.WriteString(`<h1>Metrics</h1><table><tr><td>#</td><td>Type</td><td>Name</td><td>Value</td></tr><tbody>`)

	fmt.Println(len(mtr))
	for i, metric := range mtr {
		result.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", i, metric.Type, metric.Name, metric.Value))
	}
	w.Header().Set("Content-Type", "text/html")

	io.WriteString(w, result.String())
	w.WriteHeader(http.StatusOK)
}
