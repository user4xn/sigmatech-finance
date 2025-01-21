package consumer

import (
	"clean-arch/internal/dto"
	"clean-arch/internal/factory"
	"clean-arch/pkg/tracer"
	"clean-arch/pkg/util"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

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

func (h *handler) MyProfile(c *gin.Context) {
	res, err := h.service.MyProfile(c)
	if err != nil {
		response := util.APIResponse("Failed to get my profile", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Successfully get my profile", http.StatusOK, "success", res)
	tracer.Log(c, "info", "My Profile")
	c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateMyProfile(c *gin.Context) {
	var req dto.PayloadMyProfile
	if err := c.ShouldBind(&req); err != nil {
		response := util.APIResponse("Invalid request", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := validation.ValidateStruct(&req); err != nil {
		response := util.APIResponse("Validation failed", http.StatusBadRequest, "error", err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var (
		uploadedFileKtp    string
		uploadedFileSelfie string
	)

	if req.FileKTP != nil {
		fileExt := filepath.Ext(req.FileKTP.Filename)
		allowedExt := []string{".jpg", ".png", ".jpeg", ".bmp"}
		if !util.InArrayStr(allowedExt, fileExt) {
			response := util.APIResponse(fmt.Sprintf("File ext not valid, allowed ext is %v", allowedExt), http.StatusUnprocessableEntity, "failed", nil)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		if req.FileKTP.Size > (2 * 1024 * 1024) {
			response := util.APIResponse("File size to large, max size allowed is 2Mb", http.StatusUnprocessableEntity, "failed", nil)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		fileURLKTP, err := util.SaveFile(req.FileKTP)
		if err != nil {
			response := util.APIResponse("Failed to upload file", http.StatusInternalServerError, "error", err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		appUrl := util.GetEnv("APP_URL", "http://localhost")
		appPort := util.GetEnv("APP_PORT", "8080")

		baseURL := fmt.Sprintf("%s:%s/", appUrl, appPort)
		sanitizedLink := strings.Replace(fileURLKTP, baseURL, "", 1)

		uploadedFileKtp = sanitizedLink
		req.KTPURL = fileURLKTP
	}

	if req.FileSelfie != nil {
		fileExt := filepath.Ext(req.FileSelfie.Filename)
		allowedExt := []string{".jpg", ".png", ".jpeg", ".bmp"}
		if !util.InArrayStr(allowedExt, fileExt) {
			response := util.APIResponse(fmt.Sprintf("File ext not valid, allowed ext is %v", allowedExt), http.StatusUnprocessableEntity, "failed", nil)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		if req.FileSelfie.Size > (2 * 1024 * 1024) {
			response := util.APIResponse("File size to large, max size allowed is 2Mb", http.StatusUnprocessableEntity, "failed", nil)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
		}

		fileURLSelfie, err := util.SaveFile(req.FileSelfie)
		if err != nil {
			response := util.APIResponse("Failed to upload file", http.StatusInternalServerError, "error", err.Error())
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		appUrl := util.GetEnv("APP_URL", "http://localhost")
		appPort := util.GetEnv("APP_PORT", "8080")

		baseURL := fmt.Sprintf("%s:%s/", appUrl, appPort)
		sanitizedLink := strings.Replace(fileURLSelfie, baseURL, "", 1)

		uploadedFileSelfie = sanitizedLink
		req.SelfieURL = fileURLSelfie
	}

	if err := h.service.UpdateMyProfile(c, req); err != nil {
		if req.FileKTP != nil {
			_ = util.DeleteFile(uploadedFileKtp)
		}

		if req.FileSelfie != nil {
			_ = util.DeleteFile(uploadedFileSelfie)
		}
		response := util.APIResponse("Failed to update my profile", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully update my profile", http.StatusOK, "success", nil)
	tracer.Log(c, "info", "Update My Profile")
	c.JSON(http.StatusOK, response)
}
