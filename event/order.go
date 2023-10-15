package event

import (
	"encoding/json"

	"github.com/google/uuid"
)

var (
	OrderTopic string = "order"
)

type OrderStatus int

const (
	OrderReceived OrderStatus = iota
	OrderProcessed
)

type Order struct {
	OrderID string      `json:"order_id"`
	Amount  int         `json:"amount"`
	Status  OrderStatus `json:"status"`
}

func NewOrder(amount int) Order {
	return Order{
		OrderID: uuid.New().String(),
		Amount:  amount,
		Status:  OrderReceived,
	}
}

func (o *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
