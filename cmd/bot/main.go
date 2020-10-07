package main

import (
	"bot/pkg/basics"
	"context"
	"log"
)

var (
	// token from @BotFather
	BotToken = "1385866702:AAFvpsPzErz-jF5APKJGqXcSsvimKXAXbXo"

	// heroku url
	//WebURL = "https://clothes-bot.herokuapp.com/"
)

func main() {
	bot := basics.Bot{
		BotToken: BotToken,
		//WebURL:   WebURL,
	}
	err := bot.StartBot(context.Background())
	if err != nil {
		log.Fatalf(err.Error())
	}
}
