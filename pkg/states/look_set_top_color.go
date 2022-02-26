package states

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type LookSetTopColorState struct{
	BaseState
}

func NewLookSetTopColorState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &LookSetTopColorState{BaseState: NewBase(bot, db, s3)}
}

func (s LookSetTopColorState) Do(message *tb.Message) string {
	s.db.SetTopColor(message.Sender.ID, utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0])))
	s.bot.Send(message.Sender, constants.WhatColor(strings.ToLower(constants.Bottom)), constants.ColorButtons(true))
	return LookSetBottomColorState{}.GetName()
}

func (s LookSetTopColorState) GetName() string {
	return "lookSetTopColor"
}
