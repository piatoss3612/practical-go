package main

import (
	"math/rand"
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

	for i := 0; i < 100; i++ {
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

			SleepingBarber(barbers, rand.Intn(10)+1, rand.Intn(300)+200)
		}()
	}

	wg.Wait()
}
