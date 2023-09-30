package main

import (
	"fmt"

	"github.com/fatih/color"
)

// 손님 타입 (문자열 커스텀 타입)
type Customer string

// 손님 생성자
func NewCustomer(name string) *Customer {
	customer := Customer(name)
	return &customer
}

// 손님의 이름을 문자열로 반환합니다. (Stringer 인터페이스 구현)
func (c Customer) String() string {
	return string(c)
}

// 손님이 바버샵에 들어갑니다.
func (c *Customer) EnterBarberShop(shop *BarberShop) {
	color.Green("%s(이)가 바버샵에 도착했습니다.\n", c)
	shop.ServeCustomer(c) // 바버샵에 손님을 추가합니다.
}

// 손님이 바버샵에서 나갑니다. (이발을 받았는지 여부, 이유)
func (c *Customer) LeaveBarberShop(haircut bool, reasons ...string) {
	if !haircut {
		comment := "%s(이)가 집에 돌아갑니다."
		if len(reasons) > 0 {
			comment += fmt.Sprintf("사유: [%s]", reasons[0])
		}

		color.Red(comment, c)
		return
	}

	color.Green("%s(이)가 이발을 받고 바버샵을 나갑니다.\n", c)
}

// 손님이 이발사를 깨웁니다.
func (c *Customer) WakeBarberUp(barber *Barber) {
	color.Green("%s(이)가 %s(을)를 깨웁니다.\n", c, barber.Name)
	barber.WakeUp()
}
