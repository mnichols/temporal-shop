package server

type Config struct {
	Port            string
	GeneratedAppDir string
}

func (c Config) Prefix() string {
	return "http_server"
}
