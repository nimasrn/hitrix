package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func ServiceProviderMockStripe(mock interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.StripeService,
		Build: func(_ di.Container) (interface{}, error) {
			return mock, nil
		},
	}
}
