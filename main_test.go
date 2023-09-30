package main

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	main()
	os.Exit(m.Run())
}

func TestSleepingBarber(t *testing.T) {
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		barbers := []struct {
			name            string
			cuttingDuration time.Duration
		}{
			{"철수", 0},
			{"영희", 0},
			{"영수", 0},
			{"민수", 0},
			{"민희", 0},
			{"국봉", 0},
		}

		go func() {
			defer wg.Done()

			SleepingBarber(barbers, 200)
		}()
	}

	wg.Wait()
}
