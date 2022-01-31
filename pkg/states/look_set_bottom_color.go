package states

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type LookSetBottomColorState struct{}

func (s LookSetBottomColorState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	db.SetBottomColor(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	//msg := tb.NewMessage(message.Chat.ID, constants.WhatSeason)
	//msg.ReplyMarkup = tb.NewInlineKeyboardMarkup(constants.SeasonButtons(false)...)
	go bot.Send(message.Sender, constants.WhatSeason, constants.SeasonButtons(false))
	return LookSetSeasonState{}.GetName()
}

func (s LookSetBottomColorState) GetName() string {
	return "lookSetBottomColor"
}
