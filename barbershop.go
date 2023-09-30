package main

import (
	"sync"
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	Capacity     int            // 바버샵의 최대 수용 인원
	OpenDuration time.Duration  // 바버샵의 영업 시간
	barbers      []*Barber      // 바버샵에 있는 이발사들
	customerChan chan *Customer // 바버샵에 있는 고객들의 채널
	Open         bool           // 바버샵이 영업 중인지 여부

	wg sync.WaitGroup // 바버샵의 모든 이발사들이 퇴근할 때까지 기다리기 위한 WaitGroup
	mu sync.RWMutex   // 바버샵의 상태를 변경할 때 사용하는 뮤텍스
}

func NewBarberShop(capacity int, openDuration time.Duration) *BarberShop {
	return &BarberShop{
		Capacity:     capacity,
		OpenDuration: openDuration,
		barbers:      make([]*Barber, 0),
		customerChan: make(chan *Customer, capacity),
		Open:         false,
		wg:           sync.WaitGroup{},
		mu:           sync.RWMutex{},
	}
}

func (b *BarberShop) OpenShop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Open = true // 바버샵을 오픈합니다.

	color.Blue("공지: 바버샵이 문을 열었습니다. 영업 시간은 %s입니다.\n", b.OpenDuration)

	// goroutine을 사용하여 영업 시간이 끝나면 바버샵을 닫습니다.
	go func() {
		timer := time.NewTimer(b.OpenDuration)

		<-timer.C // 영업 시간 타이머가 끝나면 바버샵을 닫습니다.

		b.CloseShop()
	}()
}

func (b *BarberShop) CloseShop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Open = false // 바버샵을 닫습니다.

	color.Blue("공지: 영업 시간이 종료되었습니다. 대기 중인 고객들을 모두 돌려보냅니다.\n")

	// 바버샵에 있는 모든 고객들을 돌려보냅니다.
	for len(b.customerChan) > 0 {
		customer := <-b.customerChan
		customer.LeaveBarberShop(false, "바버샵의 영업 시간이 종료되었습니다.")
	}

	close(b.customerChan)

	color.Blue("공지: 모든 고객들을 돌려보냈습니다. 이발사들을 퇴근시킵니다.\n")

	// 바버샵에 있는 모든 이발사들을 퇴근시킵니다.
	for _, barber := range b.barbers {
		barber.GoodToGoHome() // 이발사들이 퇴근할 준비를 합니다.

		go func(barber *Barber) {
			<-barber.Done() // 이발사들이 퇴근할 때까지 기다립니다.
			b.wg.Done()
		}(barber)
	}
}

func (b *BarberShop) IsOpen() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Open
}

func (b *BarberShop) AddBarber(barber *Barber) <-chan *Customer {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 바버샵이 닫혀있으면 이발사를 추가할 수 없습니다.
	if !b.Open {
		return nil
	}

	// 이발사를 추가하고 고객들을 받아들일 채널을 반환합니다.
	b.barbers = append(b.barbers, barber)

	b.wg.Add(1)

	return b.customerChan
}

func (b *BarberShop) ServeCustomer(customer *Customer) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 바버샵이 닫혀있으면 고객을 추가할 수 없습니다.
	if !b.Open {
		customer.LeaveBarberShop(false, "바버샵이 문을 닫았습니다.")
		return
	}

	select {
	case b.customerChan <- customer: // 바버샵에 고객을 추가합니다.
	default: // 바버샵이 꽉 찼으면 고객을 추가할 수 없습니다.
		customer.LeaveBarberShop(false, "바버샵이 꽉 찼습니다.")
	}
}

func (b *BarberShop) WaitTilAllDone() {
	b.wg.Wait() // 바버샵의 모든 이발사들이 퇴근할 때까지 기다립니다.
	color.Blue("공지: 모든 이발사가 퇴근했습니다. 바버샵이 문을 닫습니다. 다음에 또 오세요!\n")
}
