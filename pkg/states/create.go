package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
)

type CreateState struct{}

func (s CreateState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	recent, _ := strconv.Atoi(message.Text)
	bot.Send(message.Sender, constants.Added(db.GetName(message.Sender.ID, recent), recent))
	return MainState{}.GetName()
}

func (s CreateState) GetName() string {
	return "waitColor"
}

