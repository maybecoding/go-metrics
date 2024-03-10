package entity

//easyjson:json
type Metrics struct {
	ID    string   `json:"id" db:"name"`               // Имя метрики
	MType string   `json:"type" db:"type"`             // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" db:"delta"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" db:"value"` // Значение метрики в случае передачи gauge
}

//easyjson:json
type MetricsList []Metrics
