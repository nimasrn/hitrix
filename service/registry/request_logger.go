package registry

import (
	"errors"

	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
)

func ServiceProviderRequestLogger() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.RequestLoggerService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.RequestLoggerEntity") == nil {
				return nil, errors.New("you should register RequestLoggerEntity")
			}

			return requestlogger.NewDBLogger(ctn.Get(service.ClockService).(clock.IClock)), nil
		},
	}
}
