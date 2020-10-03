package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
)

type SeasonState struct{}

func (s SeasonState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.SendMeColor)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.ColorButtons(false)...)
	bot.Send(msg)
	return ColorState{}.GetName()
}

func (s SeasonState) GetName() string {
	return "waitSeason"
}
