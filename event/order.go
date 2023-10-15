package event

import (
	"encoding/json"

	"github.com/google/uuid"
)

type EventTopic string

var (
	OrderTopic EventTopic = "order"
)

type Order struct {
	OrderID string `json:"order_id"`
	Amount  int    `json:"amount"`
}

func (o *Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

func NewUUID() string {
	return uuid.New().String()
}
