package ginHandler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-payments-api/config"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

// @Summary API Health Liveness
// @Description Check API Health Liveness with some app data
// @Tags API
// @Accept json
// @Produce json
// @Router /liveness [get]
// @Success 200 {object} port.APIhealthResponse
func Liveness(c *gin.Context) {
	cfg := c.MustGet("cfg").(config.API)

	sumaryData := fmt.Sprintf("%s:%s in TagVersion: %s on Envoriment:%s responds OK",
		cfg.Name,
		cfg.Port,
		cfg.TagVersion,
		cfg.Env)
	c.JSON(http.StatusOK, &port.APIhealthResponse{
		Message: "OK",
		Sumary:  sumaryData,
	})
}
