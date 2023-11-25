package httpsender

import (
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"net/http"
	"strconv"
)

// HTTPSender Структура по отправке
type HTTPSender struct {
	address  string
	method   string
	template string
}

func (s *HTTPSender) Send(metrics []*app.Metrics) {
	for _, metric := range metrics {
		value := ""
		if metric.MType == app.MetricGauge && metric.Value != nil {
			value = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
		} else if metric.MType == app.MetricCounter && metric.Delta != nil {
			value = strconv.FormatInt(*metric.Delta, 10)
		}

		url := fmt.Sprintf(s.template, s.address, metric.MType, metric.ID, value)

		req, err := http.NewRequest(s.method, url, nil)
		if err != nil {
			fmt.Println("error due request creation: ", err)
			continue
		}
		req.Header.Add("Content-Type", "text/plan")
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println("error due sending request: ", err)
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("data hasn't sent status code is:%d \n", err)
			continue
		}
		_ = resp.Body.Close()
	}
}

func New(address string, method string, template string) *HTTPSender {
	return &HTTPSender{address: address, method: method, template: template}
}
