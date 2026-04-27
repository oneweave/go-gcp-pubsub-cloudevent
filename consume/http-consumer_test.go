package consume

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPayload struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

func TestDecodeHTTPCloudEventJSON(t *testing.T) {
	req := httptestRequest(t, `{
		"specversion":"1.0",
		"id":"evt-1",
		"source":"oneweave://orders",
		"type":"order.ready",
		"datacontenttype":"application/json",
		"time":"2026-02-17T10:11:12Z",
		"data":{"orderId":"o-7","status":"done"}
	}`)
	req.Header.Set("Content-Type", "application/cloudevents+json")

	event, payload, err := DecodeHTTPCloudEventJSON[testPayload](req)
	require.NoError(t, err)

	assert.Equal(t, "evt-1", event.ID())
	assert.Equal(t, "oneweave://orders", event.Source())
	assert.Equal(t, "order.ready", event.Type())
	assert.Equal(t, "o-7", payload.OrderID)
	assert.Equal(t, "done", payload.Status)
}

func TestDecodeHTTPCloudEventJSONErrors(t *testing.T) {
	t.Run("request required", func(t *testing.T) {
		event, payload, err := DecodeHTTPCloudEventJSON[testPayload](nil)
		require.Error(t, err)
		assert.Equal(t, "request is required", err.Error())
		assert.Empty(t, event)
		assert.Equal(t, testPayload{}, payload)
	})

	t.Run("invalid cloudevent request", func(t *testing.T) {
		req := httptestRequest(t, "{")
		event, payload, err := DecodeHTTPCloudEventJSON[testPayload](req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode cloudevent request")
		assert.Empty(t, event)
		assert.Equal(t, testPayload{}, payload)
	})

	t.Run("decode payload fails", func(t *testing.T) {
		req := httptestRequest(t, `{
			"specversion":"1.0",
			"id":"evt-2",
			"source":"oneweave://orders",
			"type":"order.ready",
			"datacontenttype":"application/json",
			"data":"not-an-object"
		}`)
		req.Header.Set("Content-Type", "application/cloudevents+json")

		event, payload, err := DecodeHTTPCloudEventJSON[testPayload](req)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode event data")
		assert.Empty(t, event)
		assert.Equal(t, testPayload{}, payload)
	})
}

func TestNewHTTPConsumer(t *testing.T) {
	httpConsumer := NewHTTPConsumer()
	require.NotNil(t, httpConsumer)
}

func TestHTTPConsumerConsumeHTTPRequest(t *testing.T) {
	t.Run("request required", func(t *testing.T) {
		httpConsumer := NewHTTPConsumer()
		event, consumeErr := httpConsumer.ConsumeHTTPRequest(nil)
		require.Error(t, consumeErr)
		assert.Equal(t, "request is required", consumeErr.Error())
		assert.Empty(t, event)
	})

	t.Run("invalid cloudevent request", func(t *testing.T) {
		httpConsumer := NewHTTPConsumer()
		req := httptestRequest(t, "{")
		event, consumeErr := httpConsumer.ConsumeHTTPRequest(req)
		require.Error(t, consumeErr)
		assert.Contains(t, consumeErr.Error(), "decode cloudevent request")
		assert.Empty(t, event)
	})

	t.Run("success", func(t *testing.T) {
		httpConsumer := NewHTTPConsumer()
		req := httptestRequest(t, `{
			"specversion":"1.0",
			"id":"evt-4",
			"source":"oneweave://orders",
			"type":"order.ready",
			"datacontenttype":"application/json",
			"data":{"orderId":"o-9","status":"done"}
		}`)
		req.Header.Set("Content-Type", "application/cloudevents+json")

		event, consumeErr := httpConsumer.ConsumeHTTPRequest(req)
		require.NoError(t, consumeErr)
		assert.Equal(t, "evt-4", event.ID())
	})
}

func httptestRequest(t *testing.T, body string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "http://example.com/push", strings.NewReader(body))
	require.NoError(t, err)
	return req
}
