package auth

import (
	"esensi-test/internal/dto"
	"esensi-test/internal/factory"
	"esensi-test/pkg/util"
	"io"
	"net/http"
	"regexp"

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

func (h *handler) Login(c *gin.Context) {
	var payload dto.PayloadLogin
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}

		response := util.APIResponse("Error login user", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	createToken, err := h.service.Login(c, payload)
	if err != nil {
		response := util.APIResponse("Error Create JWT Token", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := util.APIResponse("Success login user", http.StatusOK, "success", createToken)
	c.JSON(http.StatusOK, response)
	return

}

func (h *handler) Logout(c *gin.Context) {
	header := c.Request.Header["Authorization"]
	rep := regexp.MustCompile(`(Bearer)\s?`)
	bearerStr := rep.ReplaceAllString(header[0], "")
	err := h.service.Logout(c, bearerStr)

	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success logout user", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
