package main

import (
	"fmt"
)

type Post struct {
	Title    string
	URL      string
	Contents string
}

func main() {
	posts, done, errs := ScrapeTokenPost(true)

	for {
		select {
		case post := <-posts:
			if post == nil {
				continue
			}
			fmt.Println(post)
		case err := <-errs:
			if err == nil {
				continue
			}
			fmt.Println(err)
		case <-done:
			fmt.Println("Done")
			return
		}
	}
}
