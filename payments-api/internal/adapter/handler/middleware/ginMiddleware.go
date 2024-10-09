package ginMiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"
)

func AppInject(app bootstrap.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app", app)
		c.Next()
	}
}
