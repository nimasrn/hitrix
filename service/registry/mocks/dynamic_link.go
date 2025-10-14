package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockDynamicLink(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.DynamicLinkService,
		Build: func(_ di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
