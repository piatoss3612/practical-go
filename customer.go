package main

import "github.com/fatih/color"

type Customer string

func NewCustomer(name string) *Customer {
	customer := Customer(name)
	return &customer
}

func (c Customer) String() string {
	return string(c)
}

func (c *Customer) EnterBarberShop(shop *BarberShop) {
	color.Green("%s(이)가 바버샵에 도착했습니다.\n", c)
	shop.ServeCustomer(c) // 바버샵에 고객을 추가합니다.
}

func (c *Customer) LeaveBarberShop(haircut bool, reasons ...string) {
	if !haircut {
		reason := "%s(이)가 집에 돌아갑니다."
		if len(reasons) > 0 {
			reason += " 사유: %s"
		}

		color.Red(reason, c, reasons)
		return
	}

	color.Green("%s(이)가 이발을 받고 바버샵을 나갑니다.\n", c)
}
