package mocks

import (
	"github.com/latolukasz/fluxaorm"
	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/sms"
)

type FakeSMSSender struct {
	mock.Mock
}

func (f *FakeSMSSender) SendMessage(_ fluxaorm.Context, message *sms.Message) error {
	return f.Called(message).Error(0)
}
