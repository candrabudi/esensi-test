package auth

import "github.com/gin-gonic/gin"

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
}
