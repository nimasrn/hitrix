package apilogger

import (
	"reflect"
	"time"

	"github.com/latolukasz/fluxaorm"
)

type mysqlDBLog struct {
	logEntity  ILogEntity
	currentLog ILogEntity
}

func NewMysqlAPILogger(entity ILogEntity) IAPILogger {
	return &mysqlDBLog{logEntity: entity}
}

func (l *mysqlDBLog) LogStart(ormService fluxaorm.Context, logType string, request interface{}) {
	var logEntity ILogEntity

	if l.logEntity.GetID() == 0 {
		logEntity = l.logEntity
	} else {
		logEntity = reflect.New(reflect.ValueOf(l.logEntity).Elem().Type()).Interface().(ILogEntity)

		ormService.NewEntity(logEntity)
	}

	logEntity.SetType(logType)
	logEntity.SetRequest(request)
	logEntity.SetStatus("new")
	logEntity.SetCreatedAt(time.Now())

	ormService.Flush()

	l.currentLog = logEntity
}

func (l *mysqlDBLog) LogError(ormService fluxaorm.Context, message string, response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog

	ormService.EditEntity(currentLog)

	currentLog.SetMessage(message)
	currentLog.SetResponse(response)
	currentLog.SetStatus("failed")

	ormService.Flush()
}

func (l *mysqlDBLog) LogSuccess(ormService fluxaorm.Context, response interface{}) {
	if l.currentLog == nil {
		panic("log is not created")
	}

	currentLog := l.currentLog

	ormService.EditEntity(currentLog)

	currentLog.SetStatus("completed")
	currentLog.SetResponse(response)

	ormService.Flush()
}
