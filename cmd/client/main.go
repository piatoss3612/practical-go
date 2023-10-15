package main

import (
	"08-event-driven-kafka/event"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	c := flag.Int("c", 1, "number of clients")
	flag.Parse()

	log.Printf("Starting %d clients\n", *c)

	consumer := clientConsumer()
	defer consumer.Close()

	orderChan := make(chan event.Order)

	err := consumer.SubscribeTopics([]string{event.ClientTopic}, nil)
	if err != nil {
		panic(err)
	}

	stopChan := make(chan struct{})

	go func() {
		orders := make(map[string]event.Order)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		for {
			select {
			case <-sigChan:
				log.Println("Shutting down client...")
				close(stopChan)
				return
			case order := <-orderChan:
				_, ok := orders[order.OrderID]
				if !ok {
					orders[order.OrderID] = order
					continue
				}

				log.Printf("Order ID %s is ready for pickup. Enjoy your coffee!\n", order.OrderID)
				delete(orders, order.OrderID)
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

				_, ok := orders[order.OrderID]
				if !ok {
					orderChan <- order
					continue
				}

				log.Printf("Order ID %s is ready for pickup. Enjoy your coffee!\n", order.OrderID)
				delete(orders, order.OrderID)
			}
		}
	}()

	client := http.DefaultClient

	for i := 0; i < *c; i++ {
		go func() {
			amount := rand.Intn(5) + 1
			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:8080/order?amount=%d", amount), nil)
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

			orderChan <- order
		}()

		time.Sleep(1 * time.Second)
	}

	<-stopChan
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
