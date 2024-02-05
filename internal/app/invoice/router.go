package invoice

import (
	"esensi-test/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate())
	g.GET("/list", h.FindAll)
	g.GET("/detail/:id", h.Detail)
	g.POST("/store", h.Store)
	g.PUT("/update/:id", h.Update)
	g.PUT("/cancel-invoice/:id", h.CancelInvoice)
}
