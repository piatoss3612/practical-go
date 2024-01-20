package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	var redisURL string

	if os.Getenv("GO_ENV") == "docker" {
		redisURL = os.Getenv("DOCKER_REDIS_URL")
	} else {
		redisURL = os.Getenv("REDIS_URL")
	}

	cache, err := NewRedisCache(context.Background(), redisURL)
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

	bot, err := NewBot(os.Getenv("DISCORD_BOT_TOKEN"))
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
