package v1

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis/config"
)

func GetRedisAdditionalConfig(cfg *config.Config) (value *AdditionalConfigType) {
	value = &AdditionalConfigType{}

	if v, ok := cfg.AdditionalConfig.(*AdditionalConfigType); ok {
		value = v
	}

	return value
}
