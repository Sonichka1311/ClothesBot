package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"strings"
)

type MainState struct{}

func (s MainState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	splitMessage := strings.Split(message.Text, "_")
	switch splitMessage[0] {
	case "/upload":
		go bot.Send(tgbotapi.NewMessage(
			message.Chat.ID,
			constants.SendMePhoto,
		))
		return PhotoState{}.GetName()
	case "/wardrobe":
		Wardrobe(bot, db, message)
		return MainState{}.GetName()
	case "/thing":
		GetThing(bot, db, message)
		return MainState{}.GetName()
	case "/look":
		msg := tgbotapi.NewMessage(message.Chat.ID, constants.WhatTopColor)
		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.ColorButtons(true)...)
		go bot.Send(msg)
		return WhatTopColorState{}.GetName()
	case "/dirty":
		if len(splitMessage) > 1 {
			MakeDirty(bot, db, message)
		} else {
			Dirty(bot, db, message)
		}
		return MainState{}.GetName()
	case "/clean":
		MakeClean(bot, db, message)
		return MainState{}.GetName()
	case "/show":
		GetByType(bot, db, message)
		return MainState{}.GetName()
	case "/get":
		GetRandomThing(bot, db, message, splitMessage[1])
		return MainState{}.GetName()
	case "/delete":
		DeleteThing(bot, db, message, splitMessage[1])
		return MainState{}.GetName()
	default:
		SmthWrong(bot, message.Chat.ID)
		return MainState{}.GetName()
	}
}

func (s MainState) GetName() string {
	return "main"
}
