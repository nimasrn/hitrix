package registry

import (
	"errors"

	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/mail"
)

func ServiceProviderMail(newFunc mail.NewSenderFunc) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.MailService,
		Build: func(ctn di.Container) (interface{}, error) {
			ormEngine := ctn.Get(service.ORMEngineService).(fluxaorm.Engine)
			if ormEngine.Registry().EntitySchema("entity.MailTrackerEntity") == nil {
				return nil, errors.New("you should register MailTrackerEntity")
			}

			return mail.NewSender(
				ctn.Get(service.ConfigService).(config.IConfig),
				ctn.Get(service.ClockService).(clock.IClock),
				ctn.Get(service.ErrorLoggerService).(errorlogger.ErrorLogger),
				newFunc,
			)
		},
	}
}
