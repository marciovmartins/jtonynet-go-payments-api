package ginRoutes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-payments-api/config"
	handler "github.com/jtonynet/go-payments-api/internal/adapter/handler/http"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/handler/middleware"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"
)

func GinHandleRequests(cfg config.API, app bootstrap.App) {
	r := gin.Default()

	apiGroup := r.Group("/")
	apiGroup.Use(ginMiddleware.AppInject(app))

	apiGroup.GET("/accounts", handler.RetrieveAccountList)

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
