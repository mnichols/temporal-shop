package subscription

type Config struct {
	Port                   string
	ShowsGraphqlPlayground bool
	EncryptionKey          string
}

func (c Config) Prefix() string {
	return "sub"
}
