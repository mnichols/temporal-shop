package server

type Config struct {
	Port                   string
	GeneratedAppDir        string
	IsServingUI            bool
	ShowsGraphqlPlayground bool
	EncryptionKey          string
	SubscriptionsPort      string
}

func (c Config) Prefix() string {
	return "http_server"
}
