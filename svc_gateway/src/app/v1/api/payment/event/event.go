package event

import (
	"github.com/google/uuid"
	"github.com/sofyan48/svc_gateway/src/app/v1/api/payment/entity"
	"github.com/sofyan48/svc_gateway/src/app/v1/utility/kafka"
)

// PAYMENTEVENT ...
const PAYMENTEVENT = "payment"

// PaymentEvent ...
type PaymentEvent struct {
	Kafka kafka.KafkaLibraryInterface
}

// PaymentEventHandler ...
func PaymentEventHandler() *PaymentEvent {
	return &PaymentEvent{
		Kafka: kafka.KafkaLibraryHandler(),
	}
}

// PaymentEventInterface ...
type PaymentEventInterface interface {
	PaymentCreateEvent(data *entity.PaymentEvent) (*entity.PaymentEvent, error)
}

// PaymentCreateEvent ...
func (event *PaymentEvent) PaymentCreateEvent(data *entity.PaymentEvent) (*entity.PaymentEvent, error) {
	format := event.Kafka.GetStateFull()
	format.Action = data.Action
	format.CreatedAt = data.CreatedAt
	format.Data = data.Data
	format.UUID = uuid.New().String()
	data.UUID = format.UUID
	go event.Kafka.SendEvent(PAYMENTEVENT, format)
	data.Status = "QUEUE"
	return data, nil
}
