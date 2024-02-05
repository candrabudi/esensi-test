package invoice

import (
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/pkg/util"
	"io"
	"net/http"
	"strconv"

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
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10
	offset := 0

	if limitVal, err := strconv.Atoi(limitStr); err == nil {
		limit = limitVal
	}

	if offsetVal, err := strconv.Atoi(offsetStr); err == nil {
		offset = offsetVal
	}

	filters := map[string]interface{}{
		"invoice_no":    c.Query("invoice_no"),
		"issue_date":    c.Query("issue_date"),
		"subject":       c.Query("subject"),
		"total_item":    c.Query("total_item"),
		"customer_name": c.Query("customer_name"),
		"due_date":      c.Query("due_date"),
		"status":        c.Query("status"),
	}

	invoices, _ := h.service.FindAll(c, limit, offset, filters)

	response := util.APIResponse("Success get list invoices", http.StatusOK, "success", invoices)
	c.JSON(http.StatusOK, response)
}

func (h *handler) Store(c *gin.Context) {
	var input dto.InsertInvoice
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}

		response := util.APIResponse("Invoice store failed", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.Store(c.Request.Context(), input)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := util.APIResponse("Invoice store success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
	return
}
