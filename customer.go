package main

type Customer struct {
	Name string    // 고객의 이름
	done chan bool // 고객이 머리를 자르고 집으로 돌아갈 때까지 기다리는 채널
}

func NewCustomer(name string) *Customer {
	customer := Customer{Name: name, done: make(chan bool)}
	return &customer
}

func (c *Customer) String() string {
	return c.Name
}

func (c *Customer) EnterBarberShop(shop *BarberShop) error {
	return shop.AddCustomer(c) // 바버샵에 고객을 추가합니다.
}

func (c *Customer) LeaveBarberShop() {
	close(c.done) // 고객이 머리를 자르고 집으로 돌아갑니다.
}

func (c *Customer) Done() <-chan bool {
	return c.done // 고객이 머리를 자르고 집으로 돌아갈 때까지 기다립니다.
}
