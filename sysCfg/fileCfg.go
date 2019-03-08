package sysCfg

import (
	"go.uber.org/config"
	"go.uber.org/multierr"
)

func LoadFileCfgs(cfgFilePath, cfgKey string, cfg interface{}) error {
	var (
		errOnce, errMulti error
		provider *config.YAML
	)
	src := config.File(cfgFilePath)
	if provider, errOnce = config.NewYAML(src); errOnce != nil {
		errMulti = multierr.Append(errMulti, errOnce)
		return errMulti
	}

	if errOnce = provider.Get(cfgKey).Populate(cfg); errOnce != nil {
		errMulti = multierr.Append(errMulti, errOnce)
		return errMulti
	}

	return nil
}

