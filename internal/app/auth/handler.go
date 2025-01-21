package auth

import (
	"clean-arch/internal/dto"
	"clean-arch/internal/factory"
	"clean-arch/pkg/consts"
	"clean-arch/pkg/util"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) Logout(c *gin.Context) {
	header := c.Request.Header["Authorization"]
	rep := regexp.MustCompile(`(Bearer)\s?`)
	bearerStr := rep.ReplaceAllString(header[0], "")

	err := h.service.Logout(c, bearerStr)
	if err != nil {
		response := util.APIResponse(fmt.Sprintf("logout failed otp %s", err.Error()), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("logout successfull", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) Login(c *gin.Context) {
	var body dto.PayloadLogin
	if err := c.ShouldBind(&body); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := validation.ValidateStruct(&body,
		validation.Field(&body.Email,
			validation.Required,
		),
		validation.Field(&body.Password,
			validation.Required,
		),
	)

	if err != nil {
		response := util.APIResponse("login failed", http.StatusUnprocessableEntity, "failed", err.Error())
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data, err := h.service.LoginAttempt(c, body)
	if err == consts.UserNotFound {
		response := util.APIResponse(fmt.Sprintf("%s", consts.UserNotFound), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == consts.InvalidPassword {
		response := util.APIResponse(fmt.Sprintf("%s", consts.InvalidPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == consts.ErrorLoadLocationTime {
		response := util.APIResponse(fmt.Sprintf("%s", consts.ErrorLoadLocationTime), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == consts.ErrorGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", consts.ErrorGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == consts.EmptyGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", consts.EmptyGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Login", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
