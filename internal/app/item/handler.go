package item

import (
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/pkg/util"
	"io"
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

func (h *handler) Store(c *gin.Context) {
	var input dto.InsertItem
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}

		response := util.APIResponse("Item store failed", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.Store(c.Request.Context(), input)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := util.APIResponse("Item store success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
	return
}
