package user

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

func (h *handler) CreateUser(c *gin.Context) {
	var input dto.InsertUserRequest
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}

		response := util.APIResponse("User created failed", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		response := util.APIResponse("User created failed", http.StatusInternalServerError, "failed", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := util.APIResponse("User create success", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
	return
}
