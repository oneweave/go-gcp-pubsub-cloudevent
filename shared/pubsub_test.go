package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPubSubPushEnvelopeJSON(t *testing.T) {
	input := `{
		"deliveryAttempt": 3,
		"message": {
			"attributes": {"tenant": "dev"},
			"data": "eyJvayI6dHJ1ZX0=",
			"messageId": "mid-123",
			"publishTime": "2026-02-17T10:11:12.999999999Z"
		},
		"subscription": "projects/test/subscriptions/orders"
	}`

	var envelope PubSubPushEnvelope
	err := json.Unmarshal([]byte(input), &envelope)
	require.NoError(t, err)

	assert.Equal(t, 3, envelope.DeliveryAttempt)
	assert.Equal(t, "mid-123", envelope.Message.MessageID)
	assert.Equal(t, "2026-02-17T10:11:12.999999999Z", envelope.Message.PublishTime)
	assert.Equal(t, "dev", envelope.Message.Attributes["tenant"])
	assert.Equal(t, "projects/test/subscriptions/orders", envelope.Subscription)
}

func TestPubSubPushEnvelopeDoesNotReadLegacyFields(t *testing.T) {
	input := `{
		"message": {
			"message_id": "legacy-mid",
			"publish_time": "2026-02-17T10:11:12Z"
		}
	}`

	var envelope PubSubPushEnvelope
	err := json.Unmarshal([]byte(input), &envelope)
	require.NoError(t, err)

	assert.Empty(t, envelope.Message.MessageID)
	assert.Empty(t, envelope.Message.PublishTime)
}
