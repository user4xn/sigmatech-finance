package user

import (
	"clean-arch/internal/middleware"
	"clean-arch/pkg/consts"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate(string(consts.UserTypeAdmin)))
	g.GET("", h.FindAll)
	g.POST("/store", h.Store)
	g.GET("/:id/detail", h.FindOne)
	g.PUT("/:id/update", h.Update)
	g.DELETE("/:id/delete", h.Delete)
}
