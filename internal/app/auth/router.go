package auth

import (
	"esensi-test/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.Use(middleware.Authenticate())
	g.POST("/logout", h.Logout)
}
