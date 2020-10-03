package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"strings"
)

type WhatTopColorState struct{}

func (s WhatTopColorState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	db.SetTopColor(message.From.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.WhatBottomColor)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.ColorButtons(true)...)
	bot.Send(msg)
	return WhatBottomColorState{}.GetName()
}

func (s WhatTopColorState) GetName() string {
	return "waitTopColor"
}
