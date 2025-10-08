package main

import (
	"testing"

	"github.com/latolukasz/fluxaorm"
	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix/example/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
)

func TestApiLogger(t *testing.T) {
	createContextMyApp(t, "my-app",
		[]*service.DefinitionGlobal{
			registry.ServiceProviderAPILogger(&entity.APILogEntity{}),
		},
		nil,
	)

	apiLoggerService := service.DI().APILogger()

	ormService := service.DI().Orm()
	apiLoggerService.LogStart(ormService, entity.APILogTypeApple, nil)
	apiLoggerService.LogSuccess(ormService, nil)

	apiLoggerService.LogStart(ormService, entity.APILogTypeApple, nil)
	apiLoggerService.LogError(ormService, "Error appear", nil)

	entityIterator := fluxaorm.GetByIDs[entity.APILogEntity](ormService, 1, 2)

	assert.Len(t, entityIterator.Len(), 2)
	apiLogEntities := entityIterator.All()

	assert.Equal(t, apiLogEntities[0].Status, entity.APILogStatusCompleted)
	assert.Equal(t, apiLogEntities[1].Status, entity.APILogStatusFailed)
}
