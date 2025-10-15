package entity

import (
	"time"
)

const (
	APILogStatusNew       = "new"
	APILogStatusCompleted = "completed"
	APILogStatusFailed    = "failed"

	APILogTypeApple = "Apple"
)

type apiLogType struct {
	Apple string
}

var APILogTypeAll = apiLogType{
	Apple: APILogTypeApple,
}

type apiLogStatus struct {
	New       string
	Completed string
	Failed    string
}

var APILogStatusAll = apiLogStatus{
	New:       APILogStatusNew,
	Completed: APILogStatusCompleted,
	Failed:    APILogStatusFailed,
}

type APILogEntity struct {
	ID        uint64 `orm:"table=api_log;redisCache"`
	Type      string `orm:"enum=entity.APILogTypeAll;required"`
	Status    string `orm:"enum=entity.APILogStatusAll;required"`
	Request   string
	Response  string
	Message   string
	CreatedAt time.Time `orm:"time=true"`
}

func (e *APILogEntity) SetID(value uint64) {
	e.ID = value
}

func (e *APILogEntity) SetType(value string) {
	e.Type = value
}

func (e *APILogEntity) SetStatus(value string) {
	e.Status = value
}

func (e *APILogEntity) SetRequest(value string) {
	e.Request = value
}

func (e *APILogEntity) SetResponse(value string) {
	e.Response = value
}

func (e *APILogEntity) SetMessage(value string) {
	e.Message = value
}

func (e *APILogEntity) SetCreatedAt(value time.Time) {
	e.CreatedAt = value
}
