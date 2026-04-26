package produce

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Sender abstracts a message transport that can send a CloudEvent.
type Sender interface {
	Send(ctx context.Context, event cloudevents.Event) error
}

// ClientSender adapts a cloudevents.Client to the Sender interface.
type ClientSender struct {
	Client cloudevents.Client
}

// Send publishes a CloudEvent through the wrapped cloudevents.Client.
func (s ClientSender) Send(ctx context.Context, event cloudevents.Event) error {
	if s.Client == nil {
		return fmt.Errorf("cloudevents client is nil")
	}

	result := s.Client.Send(ctx, event)
	if cloudevents.IsUndelivered(result) {
		return result
	}

	return nil
}
