package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
)

type PhotoState struct{}

func (s PhotoState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	recent := db.GetRecent(message.From.ID)
	db.AddThing(message.From.ID, recent)
	db.SetPhoto(message.From.ID, recent, (*message.Photo)[0].FileID)

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.SendMeName))
	return NameState{}.GetName()
}

func (s PhotoState) GetName() string {
	return "waitPhoto"
}
