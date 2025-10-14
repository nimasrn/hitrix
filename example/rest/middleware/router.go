package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/coretrix/hitrix/example/rest/controller"
	"github.com/coretrix/hitrix/pkg/middleware"
)

func Router(ginEngine *gin.Engine) {
	var websocketController *controller.WebsocketController
	ginEngine.GET("/ws/", websocketController.InitConnection)

	var pingController *controller.PingController
	ginEngine.GET("/ping/", pingController.GetPingAction)

	middleware.ACLRouter(ginEngine)
}
