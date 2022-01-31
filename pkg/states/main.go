package states

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type MainState struct{}

func (s MainState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	splitMessage := strings.Split(message.Text, "_")
	switch splitMessage[0] {
	case "/upload":
		go bot.Send(message.Sender, constants.SendMePhoto)
		return UploadSetPhotoState{}.GetName()
	case "/wardrobe":
		Wardrobe(bot, db, message)
		return MainState{}.GetName()
	case "/thing":
		GetThing(bot, db, s3, message)
		return MainState{}.GetName()
	case "/look":
		go bot.Send(message.Sender, constants.WhatColor(strings.ToLower(constants.Top)), constants.ColorButtons(true))
		return LookSetTopColorState{}.GetName()
	case "/dirty":
		if len(splitMessage) > 1 {
			ChangePurity(bot, db, message, false)
		} else {
			Dirty(bot, db, message)
		}
		return MainState{}.GetName()
	case "/clean":
		ChangePurity(bot, db, message, true)
		return MainState{}.GetName()
	case "/show":
		GetByType(bot, db, message)
		return MainState{}.GetName()
	case "/get":
		GetRandomThing(bot, db, s3, message, splitMessage[1])
		return MainState{}.GetName()
	case "/delete":
		DeleteThing(bot, db, message, splitMessage[1])
		return MainState{}.GetName()
	//case "/all":
		//tb.NewMediaGroup()
		//return MainState{}.GetName()
	default:
		SmthWrong(bot, message.Sender)
		return MainState{}.GetName()
	}
}

func (s MainState) GetName() string {
	return "main"
}
