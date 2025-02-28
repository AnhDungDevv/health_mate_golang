package utils

import "github.com/gin-gonic/gin"

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
func GetRequestID(g *gin.Context) string {
	return g.GetHeader("X-Request-Id")
}
