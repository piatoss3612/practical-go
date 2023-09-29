package main

import (
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	main()           // main 함수 실행
	os.Exit(m.Run()) // 테스트 실행
}

func TestDiningPhilosophers(t *testing.T) {
	cnt := 1000 // 테스트 횟수

	oldOut := os.Stdout // 원래의 표준 출력

	_, w, _ := os.Pipe() // 표준 출력을 감추기 위해 파이프를 생성

	os.Stdout = w // 표준 출력을 파이프로 변경

	wg := sync.WaitGroup{} // 테스트를 goroutine으로 병렬 실행하기 위해 WaitGroup을 사용

	ThinkTime = 0
	EatTime = 0
	EtcTime = 0

	for i := 0; i < cnt; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			DiningPhilosophers() // 테스트할 함수 실행
		}()
	}

	wg.Wait() // 모든 테스트가 끝날 때까지 대기

	// 표준 출력을 원래대로 복구
	_ = w.Close()

	os.Stdout = oldOut
}
