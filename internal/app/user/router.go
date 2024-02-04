package user

import "github.com/gin-gonic/gin"

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/create", h.CreateUser)
}
