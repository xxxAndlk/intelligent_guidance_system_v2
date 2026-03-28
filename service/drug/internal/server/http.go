package server

import "github.com/gin-gonic/gin"

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/drugs", func(c *gin.Context) {})
	engine.GET("/drugs", func(c *gin.Context) {})
	engine.GET("/drugs/:id", func(c *gin.Context) {})
	engine.PUT("/drugs/:id", func(c *gin.Context) {})
	engine.DELETE("/drugs/:id", func(c *gin.Context) {})
	engine.PUT("/drugs/:id/stock", func(c *gin.Context) {})
	engine.POST("/drugs/:id/stock/add", func(c *gin.Context) {})
	engine.POST("/drugs/:id/stock/reduce", func(c *gin.Context) {})
	engine.GET("/drugs/:id/availability", func(c *gin.Context) {})
	engine.GET("/drugs/search", func(c *gin.Context) {})
	engine.PUT("/drugs/:id/discontinue", func(c *gin.Context) {})
	engine.PUT("/drugs/:id/activate", func(c *gin.Context) {})
}