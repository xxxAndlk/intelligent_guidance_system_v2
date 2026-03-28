package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(cfg *Config) *http.Server {
	return http.NewServer(
		http.Addr(":8006"),
	)
}

func RegisterHTTPRoutes(engine *gin.Engine) {
	engine.POST("/auth/login", func(c *gin.Context) {})
	engine.POST("/auth/logout", func(c *gin.Context) {})
	engine.POST("/auth/refresh", func(c *gin.Context) {})
	engine.POST("/auth/register", func(c *gin.Context) {})
	engine.PUT("/auth/password", func(c *gin.Context) {})
	engine.GET("/auth/permission", func(c *gin.Context) {})
	engine.GET("/users", func(c *gin.Context) {})
	engine.GET("/users/:id", func(c *gin.Context) {})
	engine.POST("/users/:id/roles", func(c *gin.Context) {})
	engine.DELETE("/users/:id/roles/:role_id", func(c *gin.Context) {})
	engine.PUT("/users/:id/activate", func(c *gin.Context) {})
	engine.PUT("/users/:id/lock", func(c *gin.Context) {})
	engine.PUT("/users/:id/unlock", func(c *gin.Context) {})
}