package config

type (
	Config struct {
		App
		Sender
	}

	Sender struct {
		Address     string
		Method      string
		Template    string
		IntervalSec int
	}
	App struct {
		CollectIntervalSec int
		SendIntervalSec    int
	}
)

func New() *Config {
	return &Config{
		App: App{
			CollectIntervalSec: 2,
			SendIntervalSec:    2,
		},

		Sender: Sender{
			Address:  "http://localhost:8080/",
			Method:   "POST",
			Template: "update/%s/%s/%s",
		},
	}
}
