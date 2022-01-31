package states

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type LookSetTopColorState struct{}

func (s LookSetTopColorState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	db.SetTopColor(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	bot.Send(message.Sender, constants.WhatColor(strings.ToLower(constants.Bottom)), constants.ColorButtons(true))
	return LookSetBottomColorState{}.GetName()
}

func (s LookSetTopColorState) GetName() string {
	return "lookSetTopColor"
}
