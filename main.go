package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

type Post struct {
	Title    string
	URL      string
	Contents string
}

type Summary struct {
	Title   string
	URL     string
	Summary string
}

func main() {
	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	posts, done, errs := ScrapeTokenPost(true)

	stop := make(chan struct{})
	summaries := make(chan *Summary)

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

			content := fmt.Sprintf("제목: %s\n\n내용: %s", post.Title, post.Contents)

			chat, err := llm.Call(ctx, []schema.ChatMessage{
				schema.SystemChatMessage{
					Content: "당신은 블록체인과 관련된 전문적인 지식을 갖추고 있습니다. 아래의 블록체인 관련 게시글을 읽고, 이를 10줄 이내로 요약해주세요.",
				},
				schema.HumanChatMessage{
					Content: content,
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			summaries <- &Summary{
				Title:   post.Title,
				URL:     post.URL,
				Summary: chat.GetContent(),
			}
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

	close(summaries)

	<-stop
}
