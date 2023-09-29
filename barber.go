package main

import (
	"errors"
	"sync"
	"time"

	"github.com/fatih/color"
)

type BarberState int8

const (
	Checking BarberState = iota
	Cutting
	Sleeping
)

type Barber struct {
	Name              string        // 이발사의 이름
	CuttingDuration   time.Duration // 이발사가 머리를 깍는데 걸리는 시간
	State             BarberState   // 이발사의 상태
	readyToGoHomeChan chan bool     // 이발사가 퇴근하기 위해 준비하는 채널
	doneChan          chan bool     // 이발사가 퇴근할 때까지 기다리는 채널

	mu sync.Mutex // 이발사의 상태를 변경할 때 사용하는 뮤텍스
}

func NewBarber(name string, cuttingDuration time.Duration) *Barber {
	return &Barber{
		Name:              name,
		CuttingDuration:   cuttingDuration,
		State:             Checking,
		readyToGoHomeChan: make(chan bool, 1),
		doneChan:          make(chan bool, 1),
		mu:                sync.Mutex{},
	}
}

func (b *Barber) GoToWork(shop *BarberShop) {
	b.mu.Lock()
	defer b.mu.Unlock()

	color.Magenta("%s(은)는 출근합니다.\n", b.Name)

	customers, err := shop.AddBarber(b) // 바버샵에 이발사를 추가합니다.
	if err != nil {
		if errors.Is(err, ErrBarberShopClosed) {
			color.Red("%s(은)는 출근하지 못했습니다. 바버샵이 문을 닫았습니다.\n", b.Name)
		}
		return
	}

	color.Magenta("%s(은)는 바버샵에서 일을 시작합니다.\n", b.Name)

	go b.acceptCustomers(customers) // 바버샵에 있는 고객들을 받아서 일을 합니다.
}

func (b *Barber) acceptCustomers(customers <-chan *Customer) {
	for {
		select {
		case <-b.readyToGoHomeChan: // 퇴근을 준비합니다.
			defer func() {
				close(b.doneChan)
				close(b.readyToGoHomeChan)
			}()
			color.Magenta("%s(은)는 퇴근을 준비합니다.\n", b.Name)
			time.Sleep(time.Millisecond * 3000)
			color.Magenta("%s(은)는 오늘 하루 일을 마치고 집으로 돌아갑니다.\n", b.Name)
			return
		case customer, ok := <-customers: // 바버샵에 있는 고객들을 받습니다.
			// 바버샵에 있는 고객들을 받았는데, 바버샵이 문을 닫았거나 고객이 없으면 다음 고객을 받습니다.
			if !ok || customer == nil {
				continue
			}

			b.mu.Lock()
			// 고객을 받았는데, 이발사가 자고 있으면 고객이 이발사를 깨우고, 이발사의 상태를 체크 중으로 변경합니다.
			if b.State == Sleeping {
				color.Green("%s(은)는 %s(을)를 깨웁니다.\n", customer, b.Name)
				b.State = Checking
			}
			b.mu.Unlock()

			b.cutHair(customer) // 이발사가 고객의 머리를 깍습니다.
		default:
			b.mu.Lock()
			// 대기 중인 고객이 없으면 이발사가 잠을 잡니다.
			if b.State == Checking {
				color.Magenta("%s(은)는 할 일이 없어 잠을 잡니다.\n", b.Name)
				b.State = Sleeping
			}
			b.mu.Unlock()
		}
	}
}

func (b *Barber) cutHair(customer *Customer) {
	color.Magenta("%s(은)는 %s의 머리를 깍습니다.\n", b.Name, customer)

	// 머리를 깍습니다.
	b.mu.Lock()
	b.State = Cutting
	b.mu.Unlock()

	time.Sleep(b.CuttingDuration)

	color.Magenta("%s(은)는 %s의 머리를 다 깍았습니다.\n", b.Name, customer)

	go customer.LeaveBarberShop() // 고객이 머리를 다 깍았으니 집으로 돌아갑니다.

	// 머리를 다 깍았으면 다음 고객을 받습니다.
	b.mu.Lock()
	b.State = Checking
	b.mu.Unlock()
}

func (b *Barber) GetReadyToGoHome() {
	b.readyToGoHomeChan <- true // 이발사가 퇴근할 준비가 되었음을 알립니다.
}

func (b *Barber) Done() <-chan bool {
	return b.doneChan // 이발사가 퇴근할 때까지 기다립니다.
}
