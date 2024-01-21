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
	defer Sync()

	cfg, err := NewConfig()
	if err != nil {
		Fatal("Failed to load config", err)
	}

	cache, err := NewRedisCache(context.Background(), func() string {
		if cfg.GoEnv == "docker" {
			return cfg.DockerRedisURL
		}
		return cfg.RedisURL
	}())
	if err != nil {
		Fatal("Failed to connect to redis cache", err)
	}

	Info("Successfully connected to redis cache")

	llm, err := openai.NewChat()
	if err != nil {
		Fatal("Failed to initialize llm chat", err)
	}

	Info("Successfully initialized llm chat")

	summarizer := NewSummarizer(llm, cache)
	tokenPostScraper := NewTokenPostScraper(true)

	summarized := make(chan *Post)

	scheduler := NewScheduler(summarizer)

	err = scheduler.AddScraper(tokenPostScraper, summarized, true)
	if err != nil {
		Fatal("Failed to add scraper", err)
	}

	Info("Successfully initialized scheduler")

	scheduler.Start()

	Info("Starting bot...")

	bot, err := NewBot(cfg.DiscordBotToken, true)
	if err != nil {
		Fatal("Failed to initialize bot", err)
	}

	err = bot.Open()
	if err != nil {
		Fatal("Failed to open bot", err)
	}

	ping := PingCommand{}

	err = bot.RegisterCommand(&ping)
	if err != nil {
		Fatal("Failed to register command", err)
	}

	Info("Successfully registered command")

	Info("Bot is running")

	<-GracefulShutdown(func() {
		Info("Gracefully shutting down bot")

		// err := bot.UnregisterCommands()
		// if err != nil {
		// 	Fatal("Failed to unregister commands", err)
		// }

		err := bot.Close()
		if err != nil {
			Fatal("Failed to close bot", err)
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
