package registry

import (
	"errors"

	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/oss"
)

func ServiceProviderOSS(newFunc oss.NewProviderFunc, namespaces oss.Namespaces) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.OSService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.OSSBucketCounterEntity") == nil {
				return nil, errors.New("you should register OSSBucketCounterEntity")
			}

			return newFunc(
				ctn.Get(service.ConfigService).(config.IConfig),
				ctn.Get(service.ClockService).(clock.IClock),
				namespaces,
			)
		},
	}
}
