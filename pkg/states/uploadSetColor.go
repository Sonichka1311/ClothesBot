package states

import (
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
)

type UploadSetColorState struct{}

func (s UploadSetColorState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	recent := db.GetRecent(message.Sender.ID)
	db.SetColor(message.Sender.ID, recent, strings.ToLower(strings.Split(message.Text, " ")[0]))
	message.Text = strconv.Itoa(recent)
	db.SetRecent(message.Sender.ID, recent+1)
	return CreateState{}.Do(bot, db, message)
}

func (s UploadSetColorState) GetName() string {
	return "uploadSetColor"
}
