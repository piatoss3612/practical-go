package main

import (
	"fmt"
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

	c.OnHTML(`div[id=content] div.list_item_title`, func(e *colly.HTMLElement) {
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

	detailCollector.OnHTML(`div[id=content] div[id=articleContentArea]`, func(e *colly.HTMLElement) {
		categories := e.ChildTexts("div.view_blockchain_item > span")
		title := e.ChildText("span.view_top_title")
		thumbnail := e.ChildAttr("div.imgBox > img", "src")

		builder := strings.Builder{}

		e.ForEach("div.article_content > p", func(_ int, h *colly.HTMLElement) {
			content := strings.TrimSpace(h.Text)
			if content == "" {
				return
			}

			if strings.Contains(content, "[email") {
				return
			}

			strongs := h.ChildTexts("strong")

			if len(strongs) > 0 {
				for _, strong := range strongs {
					content = strings.ReplaceAll(content, strong, fmt.Sprintf("**%s**", strong))
				}
			}

			builder.WriteString(fmt.Sprintf("%s\n\n", content))
		})

		posts <- &Post{
			Title:      title,
			Categories: categories,
			URL:        e.Request.URL.String(),
			Thumbnail:  thumbnail,
			Contents:   builder.String(),
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
