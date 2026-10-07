package realtime

import (
	"context"
	"time"
)

// Short lived key in the IOS application
type Secret struct {
	Value     string
	ExpiresAt time.Time
}

// The realtime session the key is created for
type SessionConfig struct {
	Model string 
	Voice string 
	Instructions string 
}

// I use dependency injection over this interface to handle the case without and API key
type Issuer interface {
	Issue(ctx context.Context, cfg SessionConfig) (Secret, error)
}
