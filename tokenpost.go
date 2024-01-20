package main

import (
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	TokenPostBaseURL  = "https://www.tokenpost.kr"
	TokenPostCacheDir = "./tokenpost_cache"
)

func ScrapeTokenPost(logging bool) (<-chan *Post, <-chan struct{}, <-chan error) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.tokenpost.kr", "tokenpost.kr"),
		colly.CacheDir(TokenPostCacheDir),
		colly.Async(),
	)

	errs := make(chan error)
	posts := make(chan *Post)

	detailCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		if logging {
			log.Println("Visiting", r.URL.String())
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if logging {
			log.Println("Error:", r.StatusCode, err)
		}
		errs <- err
	})

	c.OnHTML("div.list_item_title", func(e *colly.HTMLElement) {
		postURL := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		if postURL == "" {
			return
		}

		detailCollector.Visit(postURL)
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		if logging {
			log.Println("Visiting", r.URL.String())
		}
	})

	detailCollector.OnError(func(r *colly.Response, err error) {
		if logging {
			log.Println("Error:", r.StatusCode, err)
		}
		errs <- err
	})

	detailCollector.OnScraped(func(r *colly.Response) {
		if logging {
			log.Println("Finished", r.Request.URL.String())
		}
	})

	detailCollector.OnHTML("div#articleContentArea", func(e *colly.HTMLElement) {
		title := e.ChildText("span.view_top_title")

		builder := strings.Builder{}

		contents := e.ChildTexts("div.article_content > p")

		for i := 0; i < len(contents); i++ {
			content := strings.TrimSpace(contents[i])
			if content == "" {
				continue
			}

			if strings.Contains(content, "[email") {
				continue
			}

			builder.WriteString(contents[i])
			builder.WriteString("\n")
		}

		posts <- &Post{
			Title:    title,
			URL:      e.Request.URL.String(),
			Contents: builder.String(),
		}
	})

	done := make(chan struct{})

	go func() {
		c.Visit("https://www.tokenpost.kr/blockchain")
		c.Wait()
		detailCollector.Wait()

		close(done)
		close(posts)
		close(errs)
	}()

	return posts, done, errs
}

func RemoveTokenPostCache() error {
	return os.RemoveAll(TokenPostCacheDir)
}
