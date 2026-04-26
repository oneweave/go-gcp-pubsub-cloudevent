package produce

import (
	"context"
	"errors"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCloudEventsClient struct {
	sendResult protocol.Result
	sentEvent  cloudevents.Event
	sendCalled bool
}

func (f *fakeCloudEventsClient) Send(_ context.Context, event cloudevents.Event) protocol.Result {
	f.sendCalled = true
	f.sentEvent = event
	return f.sendResult
}

func (f *fakeCloudEventsClient) Request(_ context.Context, _ cloudevents.Event) (*cloudevents.Event, protocol.Result) {
	return nil, cloudevents.ResultACK
}

func (f *fakeCloudEventsClient) StartReceiver(_ context.Context, _ interface{}) error {
	return nil
}

func TestClientSenderSend(t *testing.T) {
	event := cloudevents.NewEvent()
	event.SetID("evt-1")
	event.SetSource("oneweave://producer")
	event.SetType("artifact.created")

	t.Run("returns error when client is nil", func(t *testing.T) {
		sender := ClientSender{}
		err := sender.Send(context.Background(), event)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cloudevents client is nil")
	})

	t.Run("returns undelivered result as error", func(t *testing.T) {
		client := &fakeCloudEventsClient{sendResult: errors.New("transport unavailable")}
		sender := ClientSender{Client: client}

		err := sender.Send(context.Background(), event)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "transport unavailable")
		assert.True(t, client.sendCalled)
		assert.Equal(t, "evt-1", client.sentEvent.ID())
	})

	t.Run("returns nil for ack result", func(t *testing.T) {
		client := &fakeCloudEventsClient{sendResult: cloudevents.ResultACK}
		sender := ClientSender{Client: client}

		err := sender.Send(context.Background(), event)
		require.NoError(t, err)
		assert.True(t, client.sendCalled)
	})

	t.Run("returns nil for nack result", func(t *testing.T) {
		client := &fakeCloudEventsClient{sendResult: cloudevents.ResultNACK}
		sender := ClientSender{Client: client}

		err := sender.Send(context.Background(), event)
		require.NoError(t, err)
		assert.True(t, client.sendCalled)
	})
}
