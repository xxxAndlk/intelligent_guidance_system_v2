package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/registrations", func(c *gin.Context) {})
	engine.GET("/registrations", func(c *gin.Context) {})
	engine.GET("/registrations/:id", func(c *gin.Context) {})
	engine.PUT("/registrations/:id/cancel", func(c *gin.Context) {})
	engine.PUT("/registrations/:id/confirm", func(c *gin.Context) {})
	engine.PUT("/registrations/:id/complete", func(c *gin.Context) {})
	engine.DELETE("/registrations/:id", func(c *gin.Context) {})
	engine.GET("/patients/:patient_id/registrations", func(c *gin.Context) {})
	engine.GET("/doctors/:doctor_id/registrations", func(c *gin.Context) {})
	engine.GET("/departments/:department_id/registrations", func(c *gin.Context) {})
}