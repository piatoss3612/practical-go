package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pandodao/tokenizer-go"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

var (
	ErrInvalidPost   = fmt.Errorf("invalid post")
	ErrTooManyTokens = fmt.Errorf("too many tokens")
)

type Summarizer struct {
	llm   *openai.Chat
	cache Cache
}

func NewSummarizer(llm *openai.Chat, cache Cache) *Summarizer {
	return &Summarizer{
		llm:   llm,
		cache: cache,
	}
}

func (s *Summarizer) Summarize(ctx context.Context, post *Post) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if post == nil {
		return ErrInvalidPost
	}

	if s.cache.Exists(ctx, post.ID) {
		var summary string
		err := s.cache.Get(ctx, post.ID, &summary)
		if err != nil {
			return err
		}

		post.Summary = summary
		post.Summarized = true

		return nil
	}

	content := post.FormatSummarizable()

	token, err := tokenizer.CalToken(content)
	if err != nil {
		return err
	}

	if token > 4500 {
		return ErrTooManyTokens
	}

	chat, err := s.llm.Call(ctx, []schema.ChatMessage{
		schema.SystemChatMessage{
			Content: "당신은 블록체인과 관련된 전문적인 지식을 갖추고 있습니다. 아래의 블록체인 관련 게시글의 본문에는 마크다운 형식(**제목**)으로 소제목이 포함되어 있을 수 있습니다. 이를 요약해주세요.",
		},
		schema.HumanChatMessage{
			Content: content,
		},
	})
	if err != nil {
		return err
	}

	summary := chat.GetContent()

	err = s.cache.Set(ctx, post.ID, summary, time.Hour*24*3)
	if err != nil {
		return err
	}

	return nil
}
