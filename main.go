package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("잠자는 이발사 문제")
	fmt.Println("==================================================")
	barbers := []struct {
		name            string
		cuttingDuration time.Duration
	}{
		{"철수", 1000 * time.Millisecond},
		{"영희", 2000 * time.Millisecond},
		{"영수", 3000 * time.Millisecond},
		{"민수", 2000 * time.Millisecond},
		{"민희", 3000 * time.Millisecond},
		{"국봉", 1000 * time.Millisecond},
	}

	SleepingBarber(barbers, 10, 400) // 잠자는 이발사 문제를 해결합니다.
	fmt.Println("==================================================")
}

func SleepingBarber(barbers []struct {
	name            string
	cuttingDuration time.Duration
}, capacity, arrivalRate int) {
	shop := NewBarberShop(capacity, time.Duration(time.Second)*10) // 10명의 고객을 수용할 수 있고 10초 동안 영업하는 바버샵을 만듭니다.

	shop.OpenShop() // 바버샵을 오픈합니다.

	// 바버샵에 이발사들을 추가합니다.
	for _, barber := range barbers {
		NewBarber(barber.name, barber.cuttingDuration).GoToWork(shop)
	}

	go randomCustomers(shop, arrivalRate) // 랜덤한 시간 간격으로 고객들이 바버샵에 들어갑니다.

	shop.WaitTilAllDone() // 바버샵이 문을 닫고 모든 이발사들이 퇴근할 때까지 기다립니다.
}

func randomCustomers(shop *BarberShop, arrivalRate int) {
	customerId := 1

	for {
		if !shop.IsOpen() {
			return
		}

		randMillisecond := rand.Int() % arrivalRate

		<-time.After(time.Duration(randMillisecond) * time.Millisecond) // 랜덤한 시간 동안 기다립니다.

		NewCustomer(fmt.Sprintf("고객%d", customerId)).EnterBarberShop(shop) // 바버샵에 고객을 추가합니다.

		customerId++ // 고객 ID를 증가시킵니다.
	}
}
