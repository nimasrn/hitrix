package mocks

import (
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
)

func FakeServiceTemplate(fake interface{}) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.HTML2PDFService,
		Build: func(_ di.Container) (interface{}, error) {
			return fake, nil
		},
	}
}
