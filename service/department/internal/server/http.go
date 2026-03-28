package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(cfg *Config) *http.Server {
	return http.NewServer(
		http.Addr(":8005"),
	)
}

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.GET("/departments", func(c *gin.Context) {})
	engine.GET("/departments/:id", func(c *gin.Context) {})
	engine.POST("/departments", func(c *gin.Context) {})
	engine.PUT("/departments/:id", func(c *gin.Context) {})
	engine.DELETE("/departments/:id", func(c *gin.Context) {})
	engine.POST("/departments/:id/doctors", func(c *gin.Context) {})
	engine.DELETE("/departments/:id/doctors/:doctor_id", func(c *gin.Context) {})
	engine.PUT("/departments/:id/duty-doctor", func(c *gin.Context) {})
	engine.PUT("/departments/:id/status", func(c *gin.Context) {})
}