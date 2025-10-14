package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func FakeServiceSTemplate(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.TemplateService,
		Build: func(_ di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
