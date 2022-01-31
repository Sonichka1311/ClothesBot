package main

import (
	"context"
	"log"

	"bot/pkg/basics"
)

var (
	// token from @BotFather
	BotToken = "<>"
)

func main() {
	//BotToken = os.Getenv("bot_token")

	bot := basics.Bot{
		BotToken: BotToken,
	}

	err := bot.StartBot(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}
}
