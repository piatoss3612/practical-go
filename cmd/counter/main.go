package main

import (
	"08-event-driven-kafka/event"
	"net/http"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Counter struct {
	*http.Server
	producer *kafka.Producer
}

func NewCounter(addr string, producer *kafka.Producer) *Counter {
	c := &Counter{
		producer: producer,
	}

	c.Server = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(c.takeOrder),
	}

	return c
}

func (c *Counter) takeOrder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.handleOrder(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Counter) handleOrder(w http.ResponseWriter, r *http.Request) {
	amount := r.URL.Query().Get("amount")
	if amount == "" {
		http.Error(w, "Missing amount", http.StatusBadRequest)
		return
	}

	numAmount, err := strconv.Atoi(amount)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	if numAmount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	order := event.NewOrder(numAmount)

	val, err := order.MarshalBinary()
	if err != nil {
		http.Error(w, "Failed to marshal order", http.StatusInternalServerError)
		return
	}

	err = c.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &event.OrderReceivedTopic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(order.OrderID),
		Value: val,
	}, nil)
	if err != nil {
		http.Error(w, "Failed to produce message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Content-Type", "application/json")
	w.Write(val)
}

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	cashier := NewCounter(":8080", producer)
	if err := cashier.ListenAndServe(); err != nil {
		panic(err)
	}
}
