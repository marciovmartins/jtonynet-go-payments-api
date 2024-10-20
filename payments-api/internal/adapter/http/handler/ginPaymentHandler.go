package ginHandler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/google/uuid"

	"github.com/jtonynet/go-payments-api/internal/bootstrap"
	"github.com/jtonynet/go-payments-api/internal/core/constants"
	"github.com/jtonynet/go-payments-api/internal/core/port"
)

// @Summary Payment Execute Transaction
// @Description Payment executes a transaction  based on the request body json data. The HTTP status is always 200. The transaction can be **approved** (code **00**), **rejected insufficient balance** (code **51**), or **rejected generally** (code **07**). [See more here](https://github.com/jtonynet/go-payments-api/tree/main?tab=readme-ov-file#about)
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body port.TransactionPaymentRequest true "Request body for Execute Transaction Payment"
// @Router /payment [post]
// @Success 200 {object} port.TransactionPaymentResponse
func PaymentExecution(ctx *gin.Context) {
	startTime := time.Now()

	defer func() {
		elapsedTime := time.Since(startTime).Milliseconds()
		fmt.Printf("Execution time: %d ms\n", elapsedTime)
	}()

	app := ctx.MustGet("app").(bootstrap.App)

	var transactionRequest port.TransactionPaymentRequest
	if err := ctx.ShouldBindBodyWith(&transactionRequest, binding.JSON); err != nil {
		fmt.Println(err)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: constants.CODE_REJECTED_GENERIC,
		})

		return
	}

	validationErrors, ok := dtoIsValid(transactionRequest)
	if !ok {
		fmt.Println(validationErrors)

		ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
			Code: constants.CODE_REJECTED_GENERIC,
		})

		return
	}

	paymentService := app.PaymentService
	returnCode, _ := paymentService.Execute(transactionRequest)

	ctx.JSON(http.StatusOK, port.TransactionPaymentResponse{
		Code: returnCode,
	})
}

func validateUUID(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(uuid.UUID)
	return ok
}

func dtoIsValid(dto any) (string, bool) {
	validate := validator.New()

	validate.RegisterValidation("uuid", validateUUID)

	err := validate.Struct(dto)
	if err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
		return strings.Join(errMsgs, ", "), false
	}

	return "", true
}
