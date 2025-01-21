package helper

import (
	"clean-arch/pkg/util"

	"github.com/gin-gonic/gin"
)

func Index(g *gin.Engine) {
	g.GET("/", func(context *gin.Context) {
		context.JSON(200, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "KKP LP Management Service",
			Version: util.GetEnv("APP_VERSION", "1.0"),
		})
	})
}
