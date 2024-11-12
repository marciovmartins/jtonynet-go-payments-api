package port

import "context"

type DBConn interface {
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	GetDB(ctx context.Context) (interface{}, error)
	GetDriver(ctx context.Context) (string, error)
}
