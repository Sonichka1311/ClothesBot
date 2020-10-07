package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type LookSetTopColorState struct{}

func (s LookSetTopColorState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	db.SetTopColor(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	bot.Send(message.Sender, constants.WhatColor(strings.ToLower(constants.Bottom)), constants.ColorButtons(true))
	return LookSetBottomColorState{}.GetName()
}

func (s LookSetTopColorState) GetName() string {
	return "lookSetTopColor"
}
