package temporal

type Config struct {
	HostPort  string
	Namespace string
}

func (c *Config) Prefix() string {
	return "temporal"
}
