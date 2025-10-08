package registry

import (
	"errors"

	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	featureflag "github.com/coretrix/hitrix/service/component/feature_flag"
)

type FeatureFlagRegistryInitFunc func(flagInterface featureflag.ServiceFeatureFlagInterface)

func ServiceProviderFeatureFlag(registry FeatureFlagRegistryInitFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.FeatureFlagEntity") == nil {
				return nil, errors.New("you should register FeatureFlagEntity")
			}

			featureFlagService := featureflag.NewFeatureFlagService(ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger))
			registry(featureFlagService)

			return featureFlagService, nil
		},
	}
}
func ServiceProviderFeatureFlagWithCache(registry FeatureFlagRegistryInitFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FeatureFlagService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.FeatureFlagEntity") == nil {
				return nil, errors.New("you should register FeatureFlagEntity")
			}

			featureFlagService := featureflag.NewFeatureFlagWithCacheService(
				ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
				ctn.Get(service.ClockService).(clock.IClock),
			)

			registry(featureFlagService)

			return featureFlagService, nil
		},
	}
}
