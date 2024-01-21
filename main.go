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
	err := SetupLogger(false, "service", "crypto gopher")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig()
	if err != nil {
		Fatal(err.Error())
	}

	cache, err := NewRedisCache(context.Background(), func() string {
		if cfg.GoEnv == "docker" {
			return cfg.DockerRedisURL
		}
		return cfg.RedisURL
	}())
	if err != nil {
		Fatal(err.Error())
	}

	Info("Successfully connected to redis cache")

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatal(err)
	}

	Info("Successfully initialized llm chat")

	_ = NewSummarizer(llm, cache)
	_ = NewTokenPostScraper(true)

	Info("Starting bot...")

	bot, err := NewBot(cfg.DiscordBotToken, true)
	if err != nil {
		log.Fatal(err)
	}

	err = bot.Run()
	if err != nil {
		log.Fatal(err)
	}

	Info("Bot is running")

	<-GracefulShutdown(func() {
		Info("Gracefully shutting down bot")

		err := bot.Close()
		if err != nil {
			log.Fatal(err)
		}
	}, syscall.SIGINT, syscall.SIGTERM)

	Info("Bot has been shutdown")
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
