package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	gocron.Scheduler
	summarizer *Summarizer
}

func NewScheduler(summarizer *Summarizer) *Scheduler {
	scheduler, _ := gocron.NewScheduler()

	return &Scheduler{
		Scheduler:  scheduler,
		summarizer: summarizer,
	}
}

func (s *Scheduler) AddScraper(scraper Scraper, res chan<- *Post, logging bool) error {
	j, err := s.NewJob(gocron.DurationJob(time.Hour*1), gocron.NewTask(func(logging bool) {
		post, done, errs := scraper.Scrape()

		for {
			select {
			case p := <-post:
				if p == nil {
					continue
				}

				if logging {
					Info("Summarizing post", p.URL)
				}

				err := s.summarizer.Summarize(context.Background(), p)
				if err != nil {
					if logging {
						Error("Failed to summarize post", err)
					}

					continue
				}

				res <- p

				fmt.Println(p.String())

				return
			case <-done:
				return
			case err := <-errs:
				if logging {
					Error("Failed to scrape post", err)
				}
			}
		}
	}, logging))
	if err != nil {
		return err
	}

	go func() {
		j.RunNow()
	}()

	return nil
}
