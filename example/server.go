package main

import (
	"log"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	model "github.com/coretrix/hitrix/example/model/socket"
	"github.com/coretrix/hitrix/pkg/middleware"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/socket"
	"github.com/coretrix/hitrix/service/registry"
)

// nolint //var eventHandlersMap is unuse
var eventHandlersMap = socket.NamespaceEventHandlerMap{
	model.DefaultWebsocketNamespace: &socket.EventHandlers{
		RegisterHandler:   model.RegisterSocketHandler,
		UnregisterHandler: model.UnRegisterSocketHandler,
	},
}

func main() {
	_, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIGlobalService(
		registry.ServiceProviderErrorLogger(),
		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrm(),
		registry.ServiceProviderClock(),
		registry.ServiceProviderJWT(),
	).RegisterDIRequestService(
		registry.ServiceProviderOrmForContext(),
	).RegisterRedisPools(
		&app.RedisPools{
			Persistent: "default",
			Cache:      "default",
			Search:     []string{"search_pool", "search_pool2"},
			Stream:     "stream_pool",
		},
	).RegisterDevPanel(&entity.DevPanelUserEntity{}, middleware.DevPanelRouter).Build()
	defer deferFunc()

	//b := &hitrix.BackgroundProcessor{Server: s}
	//b.RunAsyncOrmConsumer()
	//b.RunAsyncRequestLoggerCleaner()
	//
	//s.RunServer(9999, func(ginEngine *gin.Engine) {
	//	//middleware.RequestLogger(ginEngine, nil)
	//	exampleMiddleware.Router(ginEngine)
	//	middleware.Cors(ginEngine)
	//})
}

func auth(
	_ fluxaorm.Context,
	_ string,
	entity app.IDevPanelUserEntity,
) {
	ormEngine := service.GetServiceRequired(service.ORMEngineService).(fluxaorm.Engine)
	entitySchema := ormEngine.Registry().EntitySchema(entity)
	log.Println(1, entitySchema.GetTableName())
}
