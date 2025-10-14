package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/example/entity"
	model "github.com/coretrix/hitrix/example/model/socket"
	exampleMiddleware "github.com/coretrix/hitrix/example/rest/middleware"
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
	s, deferFunc := hitrix.New(
		"my-app", "secret",
	).RegisterDIGlobalService(

		registry.ServiceProviderConfigDirectory("config"),
		registry.ServiceProviderOrmRegistry(entity.Init),
		registry.ServiceProviderOrm(),
		registry.ServiceProviderClock(),
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

	b := &hitrix.BackgroundProcessor{Server: s}
	b.RunAsyncOrmConsumer()
	b.RunAsyncRequestLoggerCleaner()

	s.RunServer(9999, func(ginEngine *gin.Engine) {
		middleware.RequestLogger(ginEngine, nil)
		exampleMiddleware.Router(ginEngine)
		middleware.Cors(ginEngine)
	})

	e := &entity.DevPanelUserEntity{
		ID:       10,
		Username: "pass",
		Password: "pass",
	}
	auth(service.DI().Orm(), "123", e)
}

func auth(
	_ fluxaorm.Context,
	_ string,
	entity app.IDevPanelUserEntity,
) {
	ormEngine := service.GetServiceRequired(service.ORMEngineService).(fluxaorm.Engine)
	entitySchema := ormEngine.Registry().EntitySchema(entity)
	log.Println(1, entitySchema.GetTableName())
	//entitySchema.
	//q := &beeorm.RedisSearchQuery{}
	//q.FilterString(entity.GetPhoneFieldName(), phone)
	//
	//found := ormService.RedisSearchOne(entity, q)
	//found := entitySchema.
	//if !found {
	//	return "", "", errors.New("invalid credentials")
	//}
	//
	//if !entity.CanAuthenticate() {
	//	return "", "", errors.New("cannot authenticate this entity")
	//}
}
