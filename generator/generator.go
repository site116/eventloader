package generator

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofrs/uuid/v5"
	log "github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

// EventGenerator is a struct that generates events for Kafka.
type EventGenerator struct {
	batchSize int32
	topic     string
	producer  *kgo.Client
	faker     *gofakeit.Faker
	tmp       []byte
}

// NewEventGenerator creates a new generator instance.
func NewEventGenerator(
	faker *gofakeit.Faker,
	producer *kgo.Client,
	batchSize int32,
	topic string,
	tmp []byte,
) *EventGenerator {
	return &EventGenerator{
		batchSize: batchSize,
		faker:     faker,
		tmp:       tmp,
		producer:  producer,
		topic:     topic,
	}
}

// Generate generates a single event.
func (g *EventGenerator) Generate() string {
	value, err := g.faker.Template(string(g.tmp), nil)
	if err != nil {
		log.Fatal(err)
	}

	return value
}

// GenerateBatch generates a batch of events.
func (g *EventGenerator) BatchGenerate() []string {
	events := make([]string, 0, g.batchSize)
	for i := 0; i < int(g.batchSize); i++ {
		value := g.Generate()
		events = append(events, value)
	}

	return events
}

// SendEvents sends a batch of events to the broker.
func (g *EventGenerator) SendEvents(ctx context.Context, batch int) error {
	// Create a new batch
	jsonMessages := g.BatchGenerate()
	// Send the batch to the broker
	records := make([]*kgo.Record, 0, len(jsonMessages))
	for _, msg := range jsonMessages {
		records = append(records, &kgo.Record{
			Topic: g.topic,
			Value: []byte(msg),
			Key:   uuid.Must(uuid.NewV7()).Bytes(),
		})
	}

	err := g.producer.ProduceSync(ctx, records...).FirstErr()
	if err != nil {
		return err
	}
	log.Info(" Batch sent ", batch)

	return nil
}
