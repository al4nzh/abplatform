package kafka
import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"abplatform/internal/model"
)
type Producer struct {
	writer *kafka.Writer
}
func NewProducer(broker string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    "experiment-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}
func (p *Producer) SendEvent(ctx context.Context, event model.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}