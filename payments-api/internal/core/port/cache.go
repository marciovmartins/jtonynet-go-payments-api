package port

import "time"

type Cache interface {
	Readiness() error
	GetStrategy() string
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	GetDefaultExpiration() time.Duration
}
