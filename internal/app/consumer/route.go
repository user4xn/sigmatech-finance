package consumer

import (
	"clean-arch/internal/middleware"
	"clean-arch/pkg/consts"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate(string(consts.UserTypeConsumer)))
	g.GET("my-profile", h.MyProfile)
	g.PUT("my-profile/update", h.UpdateMyProfile)
}
