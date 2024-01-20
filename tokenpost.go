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

type TokenPostScraper struct {
	*colly.Collector
	logging bool
}

func NewTokenPostScraper(logging bool) *TokenPostScraper {
	return &TokenPostScraper{
		Collector: colly.NewCollector(
			colly.AllowedDomains("www.tokenpost.kr", "tokenpost.kr"),
			colly.CacheDir(TokenPostCacheDir),
			colly.Async(),
		),
		logging: logging,
	}
}

func (s *TokenPostScraper) Scrape() (<-chan *Post, <-chan struct{}, <-chan error) {
	c := s.Collector

	errs := make(chan error)
	posts := make(chan *Post)

	detailCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		if s.logging {
			log.Println("Visiting", r.URL.String())
		}
	})

	c.OnHTML(`div[id=content] div.list_item_title`, func(e *colly.HTMLElement) {
		postURL := e.Request.AbsoluteURL(e.ChildAttr("a", "href"))
		if postURL == "" {
			return
		}

		detailCollector.Visit(postURL)
	})

	c.OnScraped(func(r *colly.Response) {
		if s.logging {
			log.Println("Finished", r.Request.URL.String())
		}
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		if s.logging {
			log.Println("Visiting", r.URL.String())
		}
	})

	detailCollector.OnHTML(`div[id=content] div[id=articleContentArea]`, func(e *colly.HTMLElement) {
		categories := e.ChildTexts("div.view_blockchain_item > span")
		title := e.ChildText("span.view_top_title")
		img := e.ChildAttr("div.imgBox > img", "src")

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
			Image:      img,
			Contents:   builder.String(),
		}
	})

	detailCollector.OnScraped(func(r *colly.Response) {
		if s.logging {
			log.Println("Finished", r.Request.URL.String())
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

func (s *TokenPostScraper) ClearCache() error {
	return os.RemoveAll(TokenPostCacheDir)
}

func (s *TokenPostScraper) Close() error {
	return s.ClearCache()
}
