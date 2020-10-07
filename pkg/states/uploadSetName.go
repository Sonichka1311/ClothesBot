package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UploadSetNameState struct{}

func (s UploadSetNameState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	db.SetName(message.Sender.ID, db.GetRecent(message.Sender.ID), message.Text)
	bot.Send(message.Sender, constants.SendMeType, constants.TypeButtons())
	return UploadSetTypeState{}.GetName()
}

func (s UploadSetNameState) GetName() string {
	return "uploadSetName"
}
