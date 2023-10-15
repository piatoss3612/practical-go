package main

import (
	"08-event-driven-kafka/event"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	consumer := clientConsumer()
	defer consumer.Close()

	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/order?amount=1", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		panic("order not accepted")
	}

	var order event.Order
	err = json.NewDecoder(resp.Body).Decode(&order)
	if err != nil {
		panic(err)
	}

	log.Printf("Order accepted: %s\n", order.OrderID)

	err = consumer.SubscribeTopics([]string{event.ClientTopic}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(100 * time.Millisecond)
		if err != nil {
			continue
		}

		if msg.Key == nil || string(msg.Key) != order.OrderID {
			continue
		}

		err = order.UnmarshalBinary(msg.Value)
		if err != nil {
			log.Printf("Error unmarshalling order: %s\n", err)
			continue
		}

		log.Printf("Order status: %s\n", order.Status)
		log.Println("Enjoy the coffee!")
		return
	}
}

func clientConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "client_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	return consumer
}
