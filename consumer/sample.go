package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"
)

// ConsumerPayload consumer payload
type ConsumerPayload struct {
	ReferenceRequestID string `json:"reference_request_id,omitempty"`
}

// Ping - ping method
func (d *Handler) SampleConsumer(ctx context.Context, m *nsq.Message) error {
	payload := ConsumerPayload{}

	err := json.Unmarshal(m.Body, &payload)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%+v\n", payload)

	return nil
}
