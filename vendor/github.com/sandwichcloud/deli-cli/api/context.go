package api

import (
	"context"
	"time"
)

func CreateTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 120 * time.Second)
}
