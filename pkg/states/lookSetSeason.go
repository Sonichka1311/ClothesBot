package states

import (
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type LookSetSeasonState struct{}

func (s LookSetSeasonState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	db.SetUserSeason(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))

	return LookGenerateState{}.Do(bot, db, message)
}

func (s LookSetSeasonState) GetName() string {
	return "lookSetSeason"
}
