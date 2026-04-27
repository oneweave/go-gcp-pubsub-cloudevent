package shared

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	ID string `json:"id"`
}

func TestSetEventData(t *testing.T) {
	t.Run("sets content type and data", func(t *testing.T) {
		event := cloudevents.NewEvent(cloudevents.VersionV1)
		err := SetEventData(&event, JSONContentType, payload{ID: "a-1"})
		require.NoError(t, err)
		assert.Equal(t, JSONContentType, event.DataContentType())

		decoded, err := DecodeEventJSON[payload](event)
		require.NoError(t, err)
		assert.Equal(t, "a-1", decoded.ID)
	})

	t.Run("event required", func(t *testing.T) {
		err := SetEventData(nil, JSONContentType, payload{ID: "a-1"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "event is required")
	})

	t.Run("content type required", func(t *testing.T) {
		event := cloudevents.NewEvent(cloudevents.VersionV1)
		err := SetEventData(&event, "", payload{ID: "a-1"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "content type is required")
	})
}

func TestDecodeEventJSONErrors(t *testing.T) {
	t.Run("empty data", func(t *testing.T) {
		event := cloudevents.NewEvent(cloudevents.VersionV1)
		decoded, err := DecodeEventJSON[payload](event)
		require.Error(t, err)
		assert.Equal(t, payload{}, decoded)
		assert.Contains(t, err.Error(), "event data is empty")
	})

	t.Run("invalid json", func(t *testing.T) {
		event := cloudevents.NewEvent(cloudevents.VersionV1)
		require.NoError(t, event.SetData(JSONContentType, "not-a-json-object-for-struct"))

		decoded, err := DecodeEventJSON[payload](event)
		require.Error(t, err)
		assert.Equal(t, payload{}, decoded)
		assert.Contains(t, err.Error(), "decode event data")
	})
}
