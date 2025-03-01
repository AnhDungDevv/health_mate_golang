package utils

import (
	"context"
	"health_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ReqIDCtxKey struct{}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
func GetRequestID(c *gin.Context) string {
	return c.GetHeader("X-Request-Id")
}
func GetIPAddress(c *gin.Context) string {
	return c.Request.RemoteAddr
}

func GetRequestCtx(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), ReqIDCtxKey{}, GetRequestID(c))
}

func LogResponseError(c *gin.Context, logger logger.Logger, err error) {
	logger.Errorf(
		"ErrResponseWithLog , RequestI: %s,IPAddress: %s, Error: %s", GetRequestID(c),
		GetIPAddress(c),
		err,
	)
}

func ReadRequest(c *gin.Context, request interface{}) error {
	if err := c.ShouldBindJSON(request); err != nil {
		return err
	}
	return validate.StructCtx(context.Background(), request)
}
