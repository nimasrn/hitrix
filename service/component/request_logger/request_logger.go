package requestlogger

import (
	"net/http"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

const ID = "request_logger_id"

type IRequestLogger interface {
	LogRequest(ormService fluxaorm.Context, appName, url string, request *http.Request, contentType string) *entity.RequestLoggerEntity
	LogResponse(ormService fluxaorm.Context, requestLoggerEntity *entity.RequestLoggerEntity, responseBody []byte, status int)
}
