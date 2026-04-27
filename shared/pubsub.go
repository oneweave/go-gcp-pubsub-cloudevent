package shared

// PubSubPushEnvelope is the Google Pub/Sub push HTTP request envelope.
type PubSubPushEnvelope struct {
	DeliveryAttempt int               `json:"deliveryAttempt"`
	Message         PubSubPushMessage `json:"message"`
	Subscription    string            `json:"subscription"`
}

// PubSubPushMessage is the Pub/Sub message payload embedded in push envelopes.
type PubSubPushMessage struct {
	Attributes  map[string]string `json:"attributes"`
	Data        string            `json:"data"`
	MessageID   string            `json:"messageId"`
	PublishTime string            `json:"publishTime"`
}
