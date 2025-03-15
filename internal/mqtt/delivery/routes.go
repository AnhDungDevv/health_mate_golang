package delivery

import "github.com/gin-gonic/gin"

func NewMapMQTTRoutes(router *gin.RouterGroup, mqttHandler *MQTTHandler) {
	mqttGroup := router.Group("/mqtt")
	{
		mqttGroup.POST("/publish", mqttHandler.PublishMessage)
	}
}
