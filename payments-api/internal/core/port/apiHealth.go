package port

type APIhealthResponse struct {
	Message string `json:"message" example:"OK"`
	Sumary  string `json:"sumary" example:"payments-api:8080 in TagVersion: 0.0.0 on Envoriment:dev responds OK"`
}
