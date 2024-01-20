package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pandodao/tokenizer-go"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

type Post struct {
	Title      string   `json:"title"`
	URL        string   `json:"url"`
	Categories []string `json:"categories"`
	Image      string   `json:"image"`
	Contents   string   `json:"contents"`
	Summary    string   `json:"summary"`
	Summarized bool     `json:"summarized"`
}

func main() {
	// _, err := NewCache(context.Background(), "localhost:6379")
	// if err != nil {
	// 	panic(err)
	// }

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	scraper := NewTokenPostScraper(true)

	posts, done, errs := scraper.Scrape()

	stop := make(chan struct{})
	summaries := make(chan *Post)

	go func() {
		defer close(stop)

		for s := range summaries {
			fmt.Println("========================================")
			fmt.Println("요약된 내용:")
			fmt.Println(s.Title)
			fmt.Println(s.URL)
			fmt.Println(s.Summary)
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

				content := fmt.Sprintf("제목: %s\n\n내용:\n%s", post.Title, post.Contents)

				token, err := tokenizer.CalToken(content)
				if err != nil {
					log.Println(err)
					return
				}

				if token > 4500 {
					log.Printf("Too many tokens: %d\n", token)
					return
				}

				chat, err := llm.Call(ctx, []schema.ChatMessage{
					schema.SystemChatMessage{
						Content: "당신은 블록체인과 관련된 전문적인 지식을 갖추고 있습니다. 아래의 블록체인 관련 게시글의 본문에는 마크다운 형식으로 소제목이 포함되어 있습니다. 이를 요약해주세요.",
					},
					schema.HumanChatMessage{
						Content: content,
					},
				})
				if err != nil {
					log.Fatal(err)
				}

				post.Summary = chat.GetContent()
				post.Summarized = true

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
