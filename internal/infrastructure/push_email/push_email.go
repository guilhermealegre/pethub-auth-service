package push_email

import (
	infraConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"gopkg.in/gomail.v2"
)

func New() *gomail.Dialer {
	config, err := loadConfigs()
	if err != nil {
		panic(err)
	}

	return gomail.NewDialer(
		config.Host,
		config.Port,
		config.Email,
		config.AppPassword,
	)
}

func loadConfigs() (*config, error) {
	c := &config{}
	err := infraConfig.Load(configFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
