package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"strings"
)

type ColorState struct{}

func (s ColorState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	recent := db.GetRecent(message.From.ID)
	db.SetColor(message.From.ID, recent, strings.ToLower(strings.Split(message.Text, " ")[0]))

	bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.Added(db.GetName(message.From.ID, recent), recent)))
	db.SetRecent(message.From.ID, recent+1)
	return MainState{}.GetName()
}

func (s ColorState) GetName() string {
	return "waitColor"
}
