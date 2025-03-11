package helpers

import (
	"time"

	"golang.org/x/net/context"
)

// Helper function to create context with timeout
func CreateContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
