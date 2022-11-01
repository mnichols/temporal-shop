package temporal

type Config struct {
	HostPort     string
	Namespace    string
	CloudCertPem string
	CloudCertKey string
}

func (c *Config) Prefix() string {
	return "temporal"
}
