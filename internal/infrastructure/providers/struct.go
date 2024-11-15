package providers

type config struct {
	Providers map[string]Provider `yaml:"providers"`
}
type Provider struct {
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
	Callback string `yaml:"callback"`
}
