package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type LookSetBottomColorState struct{}

func (s LookSetBottomColorState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	db.SetBottomColor(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	//msg := tb.NewMessage(message.Chat.ID, constants.WhatSeason)
	//msg.ReplyMarkup = tb.NewInlineKeyboardMarkup(constants.SeasonButtons(false)...)
	go bot.Send(message.Sender, constants.WhatSeason, constants.SeasonButtons(false))
	return LookSetSeasonState{}.GetName()
}

func (s LookSetBottomColorState) GetName() string {
	return "lookSetBottomColor"
}
