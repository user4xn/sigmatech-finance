package http

import (
	"clean-arch/internal/app/auth"
	"clean-arch/internal/app/user"
	"clean-arch/internal/factory"
	"clean-arch/internal/middleware"
	"clean-arch/pkg/config"
	"clean-arch/pkg/helper"
	"clean-arch/pkg/tracer"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Here we define route function for user Handlers that accepts gin.Engine and factory parameters
func NewHttp(g *gin.Engine, f *factory.Factory) {
	logger, err := tracer.InitLogger(strings.ToLower(config.AppEnv()))
	if err != nil {
		panic(err)
	}
	if logger == nil {
		panic(fmt.Sprintf("Logger initialization failed %s", err))
	}

	defer logger.Sync()

	helper.Index(g)

	// Here we use cors middleware
	g.Use(middleware.CORSMiddleware())

	// Here we use logger middleware before the actual API to catch any api call from clients
	g.Use(tracer.LoggingMiddleware(logger))

	// Here we use the recovery middleware to catch a panic, if panic occurs recover the application witohut shutting it off
	g.Use(tracer.RecoverMiddleware(logger))

	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	// Here we define a router group
	v1 := g.Group("/api/v1")
	// Here we register the route from user handler
	auth.NewHandler(f).Router(v1.Group("/auth"))
	user.NewHandler(f).Router(v1.Group("/user"))
}
