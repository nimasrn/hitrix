package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/pkg/response"
	"github.com/coretrix/hitrix/service"
)

type ErrorLogController struct {
}

func (controller *ErrorLogController) GetErrors(c *gin.Context) {
	response.SuccessResponse(c, service.DI().ErrorLogger().GetErrors())
}

func (controller *ErrorLogController) GetWarnings(c *gin.Context) {
	response.SuccessResponse(c, service.DI().ErrorLogger().GetWarnings())
}

func (controller *ErrorLogController) DeleteError(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		response.ErrorResponseGlobal(c, "missing id", nil)

		return
	}

	service.DI().ErrorLogger().DeleteError(id)

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) DeleteAllErrors(c *gin.Context) {
	service.DI().ErrorLogger().DeleteAllErrors()

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) DeleteWarning(c *gin.Context) {
	id := c.Param("id")
	if len(id) <= 0 {
		response.ErrorResponseGlobal(c, "missing id", nil)

		return
	}

	service.DI().ErrorLogger().DeleteWarning(id)

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) DeleteAllWarnings(c *gin.Context) {
	service.DI().ErrorLogger().DeleteAllWarnings()

	response.SuccessResponse(c, nil)
}

func (controller *ErrorLogController) Panic(_ *gin.Context) {
	panic("Forced Panic")
}
