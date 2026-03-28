package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/files", func(c *gin.Context) {})
	engine.GET("/files", func(c *gin.Context) {})
	engine.GET("/files/:id", func(c *gin.Context) {})
	engine.DELETE("/files/:id", func(c *gin.Context) {})
}