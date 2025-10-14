package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/sentry"
)

func ServiceProviderMockSentry(mock sentry.ISentry) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.SentryService,
		Build: func(_ di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
