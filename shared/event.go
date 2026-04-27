package shared

import (
	"encoding/json"
	"fmt"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// JSONContentType is the default JSON datacontenttype for CloudEvents.
const JSONContentType = "application/json"

// DecodeEventJSON decodes a CloudEvent JSON payload into the provided generic type.
func DecodeEventJSON[T any](event cloudevents.Event) (T, error) {
	var out T
	if len(event.Data()) == 0 {
		return out, fmt.Errorf("event data is empty")
	}

	if err := json.Unmarshal(event.Data(), &out); err != nil {
		return out, fmt.Errorf("decode event data: %w", err)
	}

	return out, nil
}

// SetEventData sets CloudEvent content type and payload data.
func SetEventData(event *cloudevents.Event, contentType string, payload any) error {
	if event == nil {
		return fmt.Errorf("event is required")
	}
	if strings.TrimSpace(contentType) == "" {
		return fmt.Errorf("content type is required")
	}

	event.SetDataContentType(contentType)
	if err := event.SetData(contentType, payload); err != nil {
		return fmt.Errorf("set cloudevent data: %w", err)
	}

	return nil
}
