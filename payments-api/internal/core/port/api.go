package port

import "github.com/jtonynet/go-payments-api/internal/core/domain"

const (
	CODE_APPROVED                   = domain.CODE_APPROVED
	CODE_REJECTED_GENERIC           = domain.CODE_REJECTED_GENERIC
	CODE_REJECTED_INSUFICIENT_FUNDS = domain.CODE_REJECTED_INSUFICIENT_FUNDS
)

type TimeoutSLA int64

type APIhealthResponse struct {
	Message string `json:"message" example:"OK"`
	Sumary  string `json:"sumary" example:"payments-api:8080 in TagVersion: 0.0.0 on Envoriment:dev responds OK"`
}
