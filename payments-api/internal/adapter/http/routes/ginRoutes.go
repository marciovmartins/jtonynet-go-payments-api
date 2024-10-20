package ginRoutes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/docs"
	ginHandler "github.com/jtonynet/go-payments-api/internal/adapter/http/handler"
	ginMiddleware "github.com/jtonynet/go-payments-api/internal/adapter/http/middleware"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"
)

func GinHandleRequests(cfg config.API, app bootstrap.App) {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/"

	apiGroup := r.Group("/")
	apiGroup.Use(ginMiddleware.ConfigInject(cfg))
	apiGroup.Use(ginMiddleware.AppInject(app))

	apiGroup.GET("/liveness", ginHandler.Liveness)
	apiGroup.POST("/payment", ginHandler.PaymentExecution)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := fmt.Sprintf(":%s", cfg.Port)
	r.Run(port)
}
