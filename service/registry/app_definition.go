package registry

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

func ServiceProviderApp(app *app.App) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.AppService,
		Build: func(_ di.Container) (interface{}, error) {
			return app, nil
		},
	}
}
