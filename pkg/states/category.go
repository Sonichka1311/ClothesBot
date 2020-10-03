package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
)

type CategoryState struct{}

func (s CategoryState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	db.SetCategory(message.From.ID, db.GetRecent(message.From.ID), message.Text[1:])
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.SendMeSeason))
	return SeasonState{}.GetName()
}

func (s CategoryState) GetName() string {
	return "waitCategory"
}
