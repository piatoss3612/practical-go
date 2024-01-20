package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	cache, err := NewRedisCache(context.Background(), func() string {
		if cfg.GoEnv == "docker" {
			return cfg.DockerRedisURL
		}
		return cfg.RedisURL
	}())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Redis")

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	_ = NewSummarizer(llm, cache)
	_ = NewTokenPostScraper(true)

	bot, err := NewBot(cfg.DiscordBotToken)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running")

	<-GracefulShutdown(func() {
		err := bot.Close()
		if err != nil {
			log.Fatal(err)
		}
	}, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Gracefully shutdown")
}

func GracefulShutdown(fn func(), sig ...os.Signal) <-chan struct{} {
	stop := make(chan struct{})
	sigChan := make(chan os.Signal, 1)

	sigs := sig
	if len(sigs) == 0 {
		sigs = []os.Signal{os.Interrupt}
	}

	signal.Notify(sigChan, sigs...)

	go func() {
		<-sigChan

		signal.Stop(sigChan)

		fn()

		close(sigChan)
		close(stop)
	}()

	return stop
}
