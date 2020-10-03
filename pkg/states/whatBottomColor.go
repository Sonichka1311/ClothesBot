package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"strings"
)

type WhatBottomColorState struct{}

func (s WhatBottomColorState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	db.SetBottomColor(message.From.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.WhatSeason)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.SeasonButtons(false)...)
	go bot.Send(msg)
	return WhatSeasonState{}.GetName()
}

func (s WhatBottomColorState) GetName() string {
	return "waitBottomColor"
}
