package ocp

import "fmt"

type MomsDonaliaKiosk struct{}

func (m *MomsDonaliaKiosk) TakeOrder(o Order) {
	totalPrice := o.CalculateTotal()
	fmt.Printf("MomsDonalia: 주문 합계는 $%0.2f 입니다.\n", totalPrice)
	fmt.Println("MomsDonalia: 잠시만 기다려주세요. 주문하신 메뉴를 준비하겠습니다.")
}

type Order interface {
	CalculateTotal() float64
}

type RegularOrder struct {
	PricePerUnit float64
	Quantity     int
}

func (r RegularOrder) CalculateTotal() float64 {
	return r.PricePerUnit * float64(r.Quantity)
}

type DiscountOrder struct {
	PricePerUnit float64
	Quantity     int
	Discount     float64
}

func (o DiscountOrder) CalculateTotal() float64 {
	return (o.PricePerUnit * float64(o.Quantity)) * (1.0 - o.Discount)
}

type GiftCardOrder struct {
	PricePerUnit float64
	Quantity     int
}

func (o GiftCardOrder) CalculateTotal() float64 {
	return 0.0
}
