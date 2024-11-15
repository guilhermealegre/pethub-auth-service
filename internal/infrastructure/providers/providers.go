package providers

import (
	"fmt"
	infraConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

func New() error {

	config, err := loadConfigs()
	if err != nil {
		return err
	}

	fmt.Println(config)

	goth.UseProviders(
		google.New(config.Providers[googleStr].Key, config.Providers[googleStr].Secret, config.Providers[googleStr].Callback),
		//facebook.New(config.Providers[facebookStr].Key, config.Providers[facebookStr].Secret, config.Providers[facebookStr].Callback),
		//apple.New(config.Providers[appleStr].Key, config.Providers[appleStr].Secret, config.Providers[appleStr].Callback, nil, apple.ScopeName, apple.ScopeEmail),
	)

	return nil
}

func loadConfigs() (*config, error) {
	c := &config{}
	err := infraConfig.Load(configFile, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
