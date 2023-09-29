package main

import (
	"fmt"
	"sync"
	"time"
)

type RoundTable struct{ sync.WaitGroup } // 원형 테이블은 WaitGroup으로 구현

func (t *RoundTable) Serve(philosopher *Philosopher) {
	t.Add(1) // 철학자 한 명을 추가
	go func() {
		// 철학자가 식사를 마치면 Done()을 호출
		philosopher.Eat()
		t.Done()
	}()
}

type Fork struct{ sync.Mutex } // 포크는 뮤텍스로 구현

const HUNGER = 3 // 철학자들이 식사를 하는 횟수

var (
	ThinkTime = time.Second     // 철학자들이 생각하는 시간
	EatTime   = time.Second * 2 // 철학자들이 식사하는 시간
	EtcTime   = time.Second     // 기타 시간
)

type Philosopher struct {
	name                string // 철학자의 이름
	leftFork, rightFork *Fork  // 철학자가 사용하는 포크 (왼쪽, 오른쪽)
}

func (p *Philosopher) Eat() {
	fmt.Printf("%s가 자리에 앉음\n", p.name)
	time.Sleep(EtcTime)

	// 식사를 hunger번 반복
	for i := HUNGER; i > 0; i-- {
		// 1. 일정 시간 생각을 한다.
		fmt.Printf("%s(이)가 생각 중\n", p.name)
		time.Sleep(ThinkTime)

		// 2. 왼쪽 포크가 사용 가능해질 때까지 대기한다. 만약 사용 가능하다면 집어든다.
		p.leftFork.Lock()
		fmt.Printf("%s(이)가 왼쪽 포크를 들었음\n", p.name)

		// 3. 오른쪽 포크가 사용 가능해질 때까지 대기한다. 만약 사용 가능하다면 집어든다.
		p.rightFork.Lock()
		fmt.Printf("%s(이)가 오른쪽 포크를 들었음\n", p.name)

		// 4. 양쪽의 포크를 잡으면 일정 시간만큼 식사를 한다.
		fmt.Printf("%s(이)가 식사 중\n", p.name)
		time.Sleep(EatTime)

		// 5. 오른쪽 포크를 내려놓는다.
		p.rightFork.Unlock()
		fmt.Printf("%s(이)가 오른쪽 포크를 내려놓음\n", p.name)

		// 6. 왼쪽 포크를 내려놓는다.
		p.leftFork.Unlock()
		fmt.Printf("%s(이)가 왼쪽 포크를 내려놓음\n", p.name)

		// 7. 식사를 마치지 않았다면 1번으로 돌아간다.
		time.Sleep(EtcTime)
	}

	fmt.Println(p.name, "식사 완료!")
	time.Sleep(time.Second)

	fmt.Printf("%s(이)가 자리에서 일어남\n", p.name)
}

func main() {
	start := time.Now()
	fmt.Println("식사하는 철학자들 문제")
	fmt.Println("==================================================")

	DiningPhilosophers()

	fmt.Println("테이블 위의 모든 철학자들이 식사를 마쳤습니다.")
	fmt.Println("==================================================")
	fmt.Println("실행 시간:", time.Since(start))
}

func DiningPhilosophers() {
	table := RoundTable{
		WaitGroup: sync.WaitGroup{},
	} // 원탁을 만든다.

	names := []string{"플라톤", "아리스토텔레스", "칸트", "헤겔", "라이프니츠"} // 철학자들의 이름

	count := len(names)

	forks := make([]*Fork, count)

	// 포크를 만든다.
	for i := 0; i < count; i++ {
		forks[i] = new(Fork)
	}

	for i := 0; i < count-1; i++ {
		philosopher := Philosopher{names[i], forks[i], forks[i+1]} // i번째 철학자는 i번 포크와 i+1번 포크를 사용하며 왼쪽 포크를 먼저 집어듬

		table.Serve(&philosopher) // i번째 철학자가 식사를 시작
	}

	philosopher := Philosopher{names[count-1], forks[0], forks[count-1]} // 마지막 철학자는 오른쪽 포크를 먼저 집어듬
	table.Serve(&philosopher)                                            // 마지막 철학자가 식사를 시작

	table.Wait() // 모든 철학자들이 식사를 마칠 때까지 대기
}
