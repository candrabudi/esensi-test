package item

import (
	"esensi-test/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate())
	g.GET("/list", h.FindAll)
}
