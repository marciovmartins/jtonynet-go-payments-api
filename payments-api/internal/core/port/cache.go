package port

import "time"

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	Readiness() bool
	GetDefaultExpiration() time.Duration
}
