package main

import (
	"context"
	"github.com/aldrinleal/typebot-telegram-adapter/bothandler"
	bot "github.com/go-telegram/bot"
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Infof("Starting")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botHandler, err := bothandler.NewHandler()

	if nil != err {
		log.Fatalf("Error: %s", errorx.Decorate(err, "creating bot handler"))
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(botHandler.HandlerFunc),
		bot.WithCheckInitTimeout(30 * time.Second),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_APITOKEN"), opts...)

	if nil != err {
		log.Fatalf("Error: %s", errorx.Decorate(err, "creating bot"))
	}

	b.Start(ctx)
}
