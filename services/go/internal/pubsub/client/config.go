package client

type Config struct {
	HostPort string
	Username string
	Password string
}

func (c *Config) Prefix() string {
	return "pubsub"
}
