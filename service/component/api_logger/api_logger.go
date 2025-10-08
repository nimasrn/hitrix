package apilogger

import (
	"time"

	"github.com/latolukasz/fluxaorm"
)

type IAPILogger interface {
	LogStart(ormService fluxaorm.Context, logType string, request interface{})
	LogError(ormService fluxaorm.Context, message string, response interface{})
	LogSuccess(ormService fluxaorm.Context, response interface{})
}

type ILogEntity interface {
	beeorm.Entity
	SetID(value uint64)
	SetType(value string)
	SetStatus(value string)
	SetRequest(value interface{})
	SetResponse(value interface{})
	SetMessage(value string)
	SetCreatedAt(value time.Time)
}
