package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"sync"
)

type UploadSetTypeState struct{}

func (s UploadSetTypeState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	db.SetType(message.Sender.ID, db.GetRecent(message.Sender.ID), strings.ToLower(strings.Split(message.Text, " ")[0]))
	constants.Mutex.Lock()
	constants.MutexMap[message.Sender.ID] = &sync.Mutex{}
	constants.Mutex.Unlock()
	bot.Send(message.Sender, constants.SendMeSeason, constants.SeasonButtons(true))
	return UploadSetSeasonState{}.GetName()
}

func (s UploadSetTypeState) GetName() string {
	return "uploadSetType"
}
