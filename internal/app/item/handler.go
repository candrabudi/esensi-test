package item

import (
	"esensi-test/internal/factory"
	"esensi-test/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) FindAll(c *gin.Context) {
	search := c.Query("search")
	items, _ := h.service.FindAll(c, search)

	response := util.APIResponse("Success get list items", http.StatusOK, "success", items)
	c.JSON(http.StatusOK, response)
}
