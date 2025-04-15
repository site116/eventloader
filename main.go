package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	log "github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"gopkg.in/yaml.v2"

	"github.com/site116/eventloader/config"
	"github.com/site116/eventloader/generator"
	"github.com/site116/eventloader/worker"
)

func main() {
	ctx := context.Background()
	// Parameters command line
	templatePath := flag.String(
		"template",
		"template.json",
		"path to template file",
	)
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)

	// Read the cfg file for the values
	configBytes, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// Parse the cfg ReadFile
	var cfg config.Config
	err = yaml.Unmarshal(configBytes, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Config loaded successfully")

	// Create producer
	producer := initProducer(cfg.Broker)
	if err := producer.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	if err := producer.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to broker")

	// Read the template file for the JSON structure
	tmp, err := os.ReadFile(*templatePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Template loaded successfully")

	// Create a new faker instance
	faker := gofakeit.New(0)
	faker.Seed(0)

	// Create a new gen instance
	gen := generator.NewEventGenerator(
		faker,
		producer,
		cfg.Pool.MessagesPerBatch,
		cfg.Broker.Topic,
		tmp,
	)

	now := time.Now()
	pool := worker.NewPoll(cfg.Pool, gen.SendEvents)
	pool.Run(ctx)
	log.Info("take time:", time.Since(now))
}

// initProducer initializes a Kafka producer client with the given configuration.
func initProducer(cfg config.Broker) *kgo.Client {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
	}

	if cfg.SASL.Enabled {
		mechanism := scram.Auth{
			User: cfg.SASL.Username,
			Pass: cfg.SASL.Password,
		}.AsSha256Mechanism()

		opts = append(opts, kgo.SASL(mechanism))
	}
	producer, err := kgo.NewClient(
		opts...,
	)
	if err != nil {
		log.Fatal(err)
	}

	return producer
}
