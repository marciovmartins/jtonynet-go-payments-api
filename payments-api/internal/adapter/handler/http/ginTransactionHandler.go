package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jtonynet/go-payments-api/internal/bootstrap"
)

/*
TODO:
route for test purpose.
remove it in near future
*/
func RetrieveAccountList(ctx *gin.Context) {
	app := ctx.MustGet("app").(bootstrap.App)

	list, _ := app.AccountService.RetrieveList()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
		"sumary":  list,
	})
}
