package registry

import (
	"github.com/gin-gonic/gin"
	"github.com/latolukasz/fluxaorm"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
)

func ServiceProviderOrm() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.ORMGlobalService,
		Build: func(ctn di.Container) (interface{}, error) {
			orm := ctn.Get(service.ORMEngineService).(fluxaorm.Engine).NewContext(
				ctn.Get(service.AppService).(app.App).GlobalContext,
			)

			ormDebug, ok := ctn.Get(service.ConfigService).(config.IConfig).Bool("orm_debug")
			if ok && ormDebug {
				orm.EnableQueryDebug()
			}

			return orm, nil
		},
	}
}

func ServiceProviderOrmForContext() *service.DefinitionRequest {
	return &service.DefinitionRequest{
		Name: service.ORMRequestService,
		Build: func(c *gin.Context) (interface{}, error) {
			orm := service.GetServiceRequired(service.ORMEngineService).(fluxaorm.Engine).NewContext(c.Request.Context())

			ormDebug, ok := service.DI().Config().Bool("orm_debug")
			if ok && ormDebug {
				orm.EnableQueryDebug()
			}

			return orm, nil
		},
	}
}
