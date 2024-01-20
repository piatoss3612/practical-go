package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	// _, err := NewCache(context.Background(), "localhost:6379")
	// if err != nil {
	// 	panic(err)
	// }

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	cache := NewInMemoryCache()

	summarizer := NewSummarizer(llm, cache)

	scraper := NewTokenPostScraper(true)

	posts, done, errs := scraper.Scrape()

	stop := make(chan struct{})
	summaries := make(chan *Post)

	go func() {
		defer close(stop)

		for s := range summaries {
			fmt.Println("========================================")
			fmt.Println(s)
			fmt.Println("========================================")
		}
	}()

	cnt := 0
	ctx := context.Background()
	wg := sync.WaitGroup{}

scraper:
	for {
		select {
		case post := <-posts:
			if post == nil {
				continue
			}

			if cnt == 1 {
				continue
			}

			cnt++

			wg.Add(1)

			go func() {
				defer wg.Done()

				err := summarizer.Summarize(ctx, post)
				if err != nil {
					fmt.Println(err)
					return
				}

				summaries <- post
			}()
		case err := <-errs:
			if err == nil {
				continue
			}
			fmt.Println(err)
		case <-done:
			fmt.Println("Done")
			break scraper
		}
	}

	wg.Wait()

	close(summaries)

	<-stop
}
