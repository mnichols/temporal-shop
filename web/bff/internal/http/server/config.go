package server

type Config struct {
	Port            string
	GeneratedAppDir string
	IsServingUI     bool
}

func (c Config) Prefix() string {
	return "http_server"
}
