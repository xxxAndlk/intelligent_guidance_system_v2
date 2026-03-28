package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/notifications", func(c *gin.Context) {})
	engine.GET("/notifications", func(c *gin.Context) {})
	engine.GET("/notifications/:id", func(c *gin.Context) {})
	engine.PUT("/notifications/:id/read", func(c *gin.Context) {})
	engine.DELETE("/notifications/:id", func(c *gin.Context) {})
	engine.GET("/users/:user_id/notifications", func(c *gin.Context) {})
	engine.GET("/users/:user_id/notifications/unread", func(c *gin.Context) {})
	engine.GET("/users/:user_id/notifications/unread/count", func(c *gin.Context) {})
}