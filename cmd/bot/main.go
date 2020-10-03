package main

import (
	"bot/pkg/basics"
	"context"
	"log"
)

var (
	// token from @BotFather
	BotToken = "myAwesomeBotToken"

	// heroku url
	//WebURL = "https://project-name.herokuapp.com/"
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
