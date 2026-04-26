package produce

import (
	"context"
	"errors"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubSender struct {
	event cloudevents.Event
	err   error
}

func (s *stubSender) Send(_ context.Context, event cloudevents.Event) error {
	s.event = event
	return s.err
}

func TestPublishWrapsPayloadAsCloudEvent(t *testing.T) {
	sender := &stubSender{}
	publisher, err := NewPublisher(Config{
		Source:           "oneweave://producer",
		DefaultEventType: "artifact.ready",
		DefaultExtensions: map[string]any{
			"tenant": "dev",
		},
	}, sender)
	require.NoError(t, err)
	publisher.newID = func() string { return "evt-1" }
	publisher.now = func() time.Time { return time.Unix(100, 0).UTC() }

	payload := map[string]any{"id": "a-1"}
	got, err := publisher.Publish(context.Background(), "", payload, WithSubject("artifacts"), WithExtension("region", "us-east-1"))
	require.NoError(t, err)
	assert.Equal(t, "evt-1", got.ID())
	assert.Equal(t, "artifact.ready", got.Type())
	assert.Equal(t, "oneweave://producer", got.Source())
	assert.Equal(t, "artifacts", got.Subject())
	assert.Equal(t, "application/json", got.DataContentType())
	assert.Equal(t, "dev", got.Extensions()["tenant"])
	assert.Equal(t, "us-east-1", got.Extensions()["region"])
	assert.Equal(t, "evt-1", sender.event.ID())
}

func TestPublishToTopicSetsSubject(t *testing.T) {
	sender := &stubSender{}
	publisher, err := NewPublisher(Config{Source: "oneweave://producer", DefaultEventType: "evt"}, sender)
	require.NoError(t, err)

	event, err := publisher.PublishToTopic(context.Background(), "topic-A", "evt", map[string]string{"ok": "true"})
	require.NoError(t, err)
	assert.Equal(t, "topic-A", event.Subject())
}

func TestNewPublisherValidation(t *testing.T) {
	t.Run("sender required", func(t *testing.T) {
		publisher, err := NewPublisher(Config{Source: "oneweave://producer"}, nil)
		require.Error(t, err)
		assert.Nil(t, publisher)
		assert.Contains(t, err.Error(), "sender is required")
	})

	t.Run("source required", func(t *testing.T) {
		sender := &stubSender{}
		publisher, err := NewPublisher(Config{}, sender)
		require.Error(t, err)
		assert.Nil(t, publisher)
		assert.Contains(t, err.Error(), "source is required")
	})
}

func TestPublishValidationAndErrors(t *testing.T) {
	sender := &stubSender{}
	publisher, err := NewPublisher(Config{Source: "oneweave://producer"}, sender)
	require.NoError(t, err)

	t.Run("payload required", func(t *testing.T) {
		event, err := publisher.Publish(context.Background(), "evt", nil)
		require.Error(t, err)
		assert.Equal(t, cloudevents.Event{}, event)
		assert.Contains(t, err.Error(), "payload is required")
	})

	t.Run("event type required", func(t *testing.T) {
		event, err := publisher.Publish(context.Background(), "", map[string]string{"k": "v"})
		require.Error(t, err)
		assert.Equal(t, cloudevents.Event{}, event)
		assert.Contains(t, err.Error(), "event type is required")
	})

	t.Run("sender error wrapped", func(t *testing.T) {
		sender.err = errors.New("transport down")
		event, err := publisher.Publish(context.Background(), "evt", map[string]string{"k": "v"})
		require.Error(t, err)
		assert.Equal(t, cloudevents.Event{}, event)
		assert.Contains(t, err.Error(), "send cloudevent")
		assert.Contains(t, err.Error(), "transport down")
		sender.err = nil
	})
}

func TestPublishWithContentTypeOverride(t *testing.T) {
	sender := &stubSender{}
	publisher, err := NewPublisher(Config{Source: "oneweave://producer", DefaultEventType: "evt"}, sender)
	require.NoError(t, err)

	event, err := publisher.Publish(
		context.Background(),
		"",
		map[string]string{"ok": "true"},
		WithDataContentType("application/cloudevents+json"),
	)
	require.NoError(t, err)
	assert.Equal(t, "application/cloudevents+json", event.DataContentType())
}
