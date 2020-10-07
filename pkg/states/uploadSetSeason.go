package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UploadSetSeasonState struct{}

func (s UploadSetSeasonState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	bot.Send(message.Sender, constants.SendMeColor, constants.ColorButtons(false))
	return UploadSetColorState{}.GetName()
}

func (s UploadSetSeasonState) GetName() string {
	return "uploadSetSeason"
}
