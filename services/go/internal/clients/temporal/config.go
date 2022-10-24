package temporal

type Config struct {
	HostPort     string
	Namespace    string
	CertFilePath string `split_words:"true"`
	KeyFilePath  string `split_words:"true"`
}

func (c *Config) Prefix() string {
	return "temporal"
}
