package port

type DBConn interface {
	GetDB() interface{}
	Readiness() error
	GetStrategy() string
	GetDriver() string
}
