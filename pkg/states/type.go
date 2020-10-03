package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"strings"
	"sync"
)

type TypeState struct{}

func (s TypeState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	db.SetType(message.From.ID, db.GetRecent(message.From.ID), strings.ToLower(strings.Split(message.Text, " ")[0]))
	constants.Mutex.Lock()
	constants.MutexMap[message.From.ID] = &sync.Mutex{}
	constants.Mutex.Unlock()
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.SendMeSeason)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(constants.SeasonButtons(true)...)
	bot.Send(msg)
	return SeasonState{}.GetName()
}

func (s TypeState) GetName() string {
	return "waitType"
}
