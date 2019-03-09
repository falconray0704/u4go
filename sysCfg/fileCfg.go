package sysCfg

import (
	"go.uber.org/config"
)

func LoadFileCfgs(cfgFilePath, cfgKey string, cfg interface{}) error {
	var (
		errOnce error
		provider *config.YAML
	)

	src := config.File(cfgFilePath)
	if provider, errOnce = config.NewYAML(src); errOnce != nil {
		return errOnce
	}

	if errOnce = provider.Get(cfgKey).Populate(cfg); errOnce != nil {
		return errOnce
	}

	return nil
}

