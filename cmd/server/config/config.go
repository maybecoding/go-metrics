package config

type (
	Config struct {
		Server
	}

	Server struct {
		Address string
	}
)

func NewConfig() *Config {
	return &Config{Server{Address: ":8080"}}
}
