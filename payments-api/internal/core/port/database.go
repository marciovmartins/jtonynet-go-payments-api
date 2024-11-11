package port

import "context"

type DBConn interface {
	GetDB(ctx context.Context) (interface{}, error)
	Readiness(ctx context.Context) error
	GetStrategy(ctx context.Context) (string, error)
	GetDriver(ctx context.Context) (string, error)
}
