package server

type Config struct {
	Port            string
	GeneratedAppDir string
	IsServingUI     bool
	EncryptionKey   string
}

func (c Config) Prefix() string {
	return "http_server"
}
