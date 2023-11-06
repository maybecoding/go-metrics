package httpsender

import (
	"fmt"
	"github.com/maybecoding/go-metrics.git/internal/agent/app"
	"net/http"
)

// HTTPSender Пусть для общего развития это будет функция
type HTTPSender func(metrics []*app.Metric)

func (s HTTPSender) Send(metrics []*app.Metric) {
	s(metrics)
}

func New(address string, method string, template string) HTTPSender {
	return func(metrics []*app.Metric) {
		for _, metric := range metrics {
			url := fmt.Sprintf(template, address, metric.Type, metric.Name, metric.Value)

			req, err := http.NewRequest(method, url, nil)
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
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				fmt.Printf("data hasn't sent status code is:%d \n", err)
				continue
			}
		}
	}
}
