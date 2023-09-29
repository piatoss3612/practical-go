package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

func main() {
	fmt.Println("잠자는 이발사 문제")
	fmt.Println("==================================================")

	shop := NewBarberShop(10, time.Duration(time.Second)*10) // 10명의 고객을 수용할 수 있고 10초 동안 영업하는 바버샵을 만듭니다.

	shop.OpenShop() // 바버샵을 오픈합니다.

	// 바버샵에 이발사들을 추가합니다.
	NewBarber("철수", 1000*time.Millisecond).GoToWork(shop)
	NewBarber("영희", 2000*time.Millisecond).GoToWork(shop)
	NewBarber("영수", 3000*time.Millisecond).GoToWork(shop)
	NewBarber("민수", 2000*time.Millisecond).GoToWork(shop)
	NewBarber("민희", 3000*time.Millisecond).GoToWork(shop)
	NewBarber("국봉", 1000*time.Millisecond).GoToWork(shop)

	go randomCustomers(shop) // 랜덤한 시간 간격으로 고객들이 바버샵에 들어갑니다.

	shop.WaitTilAllDone() // 바버샵이 문을 닫고 모든 이발사들이 퇴근할 때까지 기다립니다.

	fmt.Println("==================================================")
}

func randomCustomers(shop *BarberShop) {
	customerId := 1

	for {
		randMillisecond := rand.Int() % 300

		<-time.After(time.Duration(randMillisecond) * time.Millisecond) // 랜덤한 시간 동안 기다립니다.

		c := NewCustomer(fmt.Sprintf("고객%d", customerId)) // 새로운 고객을 만듭니다.

		color.Green("%s(이)가 바버샵에 들어갑니다.\n", c)

		err := c.EnterBarberShop(shop) // 고객이 바버샵에 들어갑니다.
		if err != nil {
			if errors.Is(err, ErrBarberShopClosed) {
				color.Red("%s(이)가 바버샵에 들어가지 못했습니다. 바버샵이 문을 닫았습니다.\n", c)
				return
			} else if errors.Is(err, ErrorCustomerFull) {
				color.Red("바버샵이 꽉 찼습니다. %s(은)는 집으로 돌아갑니다.\n", c)
			} else {
				color.Red("알 수 없는 오류가 발생했습니다. %s(은)는 집으로 돌아갑니다.\n", c)
			}
		}

		customerId++ // 고객 ID를 증가시킵니다.
	}
}
