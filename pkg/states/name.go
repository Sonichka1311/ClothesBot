package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
)

type NameState struct{}

func (s NameState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	db.SetName(message.From.ID, db.GetRecent(message.From.ID), message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.SendMeType)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.TypeButtons()...)
	bot.Send(msg)
	return TypeState{}.GetName()
}

func (s NameState) GetName() string {
	return "waitName"
}
