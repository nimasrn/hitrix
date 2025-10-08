package registry

import (
	"errors"

	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/translation"
)

func ServiceProviderTranslation() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.TranslationService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.TranslationTextEntity") == nil {
				return nil, errors.New("you should register TranslationTextEntity")
			}

			return translation.NewTranslationService(ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger)), nil
		},
	}
}
