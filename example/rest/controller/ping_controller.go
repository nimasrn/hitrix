package controller

import (
	"github.com/coretrix/hitrix/pkg/response"
	"github.com/gin-gonic/gin"
)

type PingController struct {
}

func (controller *PingController) GetPingAction(c *gin.Context) {
	response.SuccessResponse(c, "Pong")
}
