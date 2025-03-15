package delivery

import (
	"health_backend/internal/mqtt/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MQTTHandler struct {
	mqttUC *usecase.MQTTUsecase
}

func NewMQTTHandler(mqttUC *usecase.MQTTUsecase) *MQTTHandler {
	return &MQTTHandler{mqttUC: mqttUC}
}

func (h *MQTTHandler) PublishMessage(c *gin.Context) {
	var request struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.mqttUC.SendMessage(request.Topic, request.Message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Published successfully"})

}
