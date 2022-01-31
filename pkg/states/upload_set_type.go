package states

import (
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type UploadSetTypeState struct{}

func (s UploadSetTypeState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	db.SetType(message.Sender.ID, db.GetUser(message.Sender.ID).LastFileID, strings.ToLower(strings.Split(message.Text, " ")[0]))

	constants.Mutex.Lock()
	constants.MutexMap[message.Sender.ID] = &sync.Mutex{}
	constants.Mutex.Unlock()

	bot.Send(message.Sender, constants.SendMeSeason, constants.SeasonButtons(true))
	return UploadSetSeasonState{}.GetName()
}

func (s UploadSetTypeState) GetName() string {
	return "uploadSetType"
}
