package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/payments", func(c *gin.Context) {})
	engine.GET("/payments", func(c *gin.Context) {})
	engine.GET("/payments/:id", func(c *gin.Context) {})
	engine.PUT("/payments/:id/process", func(c *gin.Context) {})
	engine.GET("/payments/:id/status", func(c *gin.Context) {})
	engine.POST("/payments/:id/refund", func(c *gin.Context) {})
	engine.DELETE("/payments/:id", func(c *gin.Context) {})
	engine.GET("/patients/:patient_id/payments", func(c *gin.Context) {})
}