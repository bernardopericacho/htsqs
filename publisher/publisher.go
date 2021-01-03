package publisher

import (
	"context"
	"encoding/json"
)

// Publisher is the interface clients can use to publish messages
type Publisher interface {
	Publish(ctx context.Context, msg json.Marshaler) error
}
