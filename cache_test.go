package main

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryCache(t *testing.T) {
	cache := NewInMemoryCache()

	post := &Post{
		Title:      "제목",
		Contents:   "내용",
		URL:        "https://www.google.com",
		Categories: []string{"카테고리1", "카테고리2"},
		Image:      "https://www.google.com",
		Summary:    "요약",
		Summarized: true,
	}

	err := cache.Set(context.Background(), "key", post, time.Second*3)
	if err != nil {
		t.Error(err)
	}

	var post2 Post

	ok := cache.Exists(context.Background(), "key")
	if !ok {
		t.Errorf("expected %t, got %t", true, ok)
	}

	err = cache.Get(context.Background(), "key", &post2)
	if err != nil {
		t.Error(err)
	}

	if post.Title != post2.Title {
		t.Errorf("expected %s, got %s", post.Title, post2.Title)
	}

	if post.Contents != post2.Contents {
		t.Errorf("expected %s, got %s", post.Contents, post2.Contents)
	}

	if post.URL != post2.URL {
		t.Errorf("expected %s, got %s", post.URL, post2.URL)
	}

	if post.Image != post2.Image {
		t.Errorf("expected %s, got %s", post.Image, post2.Image)
	}

	if post.Summary != post2.Summary {
		t.Errorf("expected %s, got %s", post.Summary, post2.Summary)
	}

	if post.Summarized != post2.Summarized {
		t.Errorf("expected %t, got %t", post.Summarized, post2.Summarized)
	}

	if len(post.Categories) != len(post2.Categories) {
		t.Errorf("expected %d, got %d", len(post.Categories), len(post2.Categories))
	}

	for i := range post.Categories {
		if post.Categories[i] != post2.Categories[i] {
			t.Errorf("expected %s, got %s", post.Categories[i], post2.Categories[i])
		}
	}

	time.Sleep(time.Second * 4)

	err = cache.Get(context.Background(), "key", &post2)
	if err != ErrCacheMiss {
		t.Errorf("expected %v, got %v", ErrCacheMiss, err)
	}

	err = cache.Set(context.Background(), "key", post)
	if err != nil {
		t.Error(err)
	}

	err = cache.Delete(context.Background(), "key")
	if err != nil {
		t.Error(err)
	}

	err = cache.Get(context.Background(), "key", &post2)
	if err == nil {
		t.Errorf("expected %v, got %v", ErrCacheMiss, err)
	}
}
