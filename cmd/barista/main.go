package main

import (
	"08-event-driven-kafka/event"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	consumer := orderConsumer()
	defer consumer.Close()

	producer := clientProducer()
	defer producer.Close()

	err := consumer.SubscribeTopics([]string{event.OrderTopic}, nil)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	run := true

	for run {
		select {
		case <-sigChan:
			run = false
		default:
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				continue
			}

			order := event.Order{}
			err = order.UnmarshalBinary(msg.Value)
			if err != nil {
				log.Printf("Error unmarshalling order: %s\n", err)
				continue
			}

			log.Printf("Order received: %s\n", order.OrderID)

			go func() {
				time.Sleep(5 * time.Second)

				order.Status = event.OrderProcessed

				val, err := order.MarshalBinary()
				if err != nil {
					log.Printf("Error marshalling order: %s\n", err)
					return
				}

				err = producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &order.OrderID,
						Partition: kafka.PartitionAny,
					},
					Key:   []byte(order.OrderID),
					Value: val,
				}, nil)
				if err != nil {
					log.Printf("Error producing order: %s\n", err)
					return
				}

				log.Printf("Order processed: %s\n", order.OrderID)
			}()
		}
	}
}

func orderConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     "localhost:9092",
		"group.id":              "barista_group",
		"auto.offset.reset":     "earliest",
		"broker.address.family": "v4",
	})
	if err != nil {
		panic(err)
	}

	return consumer
}

func clientProducer() *kafka.Producer {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	return producer
}
