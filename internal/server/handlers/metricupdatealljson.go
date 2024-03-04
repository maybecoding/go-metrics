package handlers

import (
	"github.com/mailru/easyjson"
	"github.com/maybecoding/go-metrics.git/internal/server/entity"
	"github.com/maybecoding/go-metrics.git/pkg/logger"
	"net/http"
)

func (c *Handler) metricUpdateAllJSON(w http.ResponseWriter, r *http.Request) {

	status := http.StatusOK
	defer func() {
		w.WriteHeader(status)
	}()

	var mts entity.MetricsList
	err := easyjson.UnmarshalFromReader(r.Body, &mts)

	//decoder := json.NewDecoder(r.Body) Оптимизировано инкремент #16
	//defer func() {
	//	_ = r.Body.Close()
	//}()
	//
	//var mts []entity.Metrics
	//err := decoder.Decode(&mts)
	if err != nil {
		status = http.StatusBadRequest
		logger.Debug().Err(err).Msg("error due decode request")
		return
	}

	err = c.metric.SetAll(mts)
	if err != nil {
		status = http.StatusInternalServerError
		logger.Debug().Err(err).Msg("error due set all")
		return
	}
}
