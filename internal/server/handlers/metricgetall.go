package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/maybecoding/go-metrics.git/internal/server/metricservice"
)

func (c *Handler) metricGetAll(w http.ResponseWriter, r *http.Request) {
	mtr, err := c.metric.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	var result strings.Builder
	result.WriteString(`<h1>Metrics</h1><table><tr><td>#</td><td>Type</td><td>Name</td><td>Value</td></tr><tbody>`)

	for i, m := range mtr {
		value := ""
		if m.MType == metricservice.Gauge {
			value = strconv.FormatFloat(*m.Value, FmtFloat, -1, 64)
		} else if m.MType == metricservice.Counter {
			value = strconv.FormatInt(*m.Delta, 10)
		}
		result.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td></tr>", i, m.MType, m.ID, value))
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, result.String())

}
