package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/surgicals", func(c *gin.Context) {})
	engine.GET("/surgicals", func(c *gin.Context) {})
	engine.GET("/surgicals/:id", func(c *gin.Context) {})
	engine.PUT("/surgicals/:id/start", func(c *gin.Context) {})
	engine.PUT("/surgicals/:id/complete-step", func(c *gin.Context) {})
	engine.PUT("/surgicals/:id/complete", func(c *gin.Context) {})
	engine.PUT("/surgicals/:id/cancel", func(c *gin.Context) {})
	engine.DELETE("/surgicals/:id", func(c *gin.Context) {})
	engine.GET("/patients/:patient_id/surgicals", func(c *gin.Context) {})
}