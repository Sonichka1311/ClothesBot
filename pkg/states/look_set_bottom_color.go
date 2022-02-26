package states

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type LookSetBottomColorState struct {
	BaseState
}

func NewLookSetBottomColorState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &LookSetBottomColorState{BaseState: NewBase(bot, db, s3)}
}

func (s LookSetBottomColorState) Do(message *tb.Message) string {
	s.db.SetBottomColor(message.Sender.ID, utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0])))
	s.bot.Send(message.Sender, constants.WhatSeason, constants.SeasonButtons(false))
	return LookSetSeasonState{}.GetName()
}

func (s LookSetBottomColorState) GetName() string {
	return "lookSetBottomColor"
}
