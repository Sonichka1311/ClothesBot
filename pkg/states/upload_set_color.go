package states

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
)

type UploadSetColorState struct{}

func (s UploadSetColorState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	recent := db.GetUser(message.Sender.ID).LastFileID
	db.SetColor(message.Sender.ID, recent, strings.ToLower(strings.Split(message.Text, " ")[0]))
	message.Text = strconv.Itoa(recent)
	db.SetRecent(message.Sender.ID, recent+1)
	return CreateState{}.Do(bot, db, s3, message)
}

func (s UploadSetColorState) GetName() string {
	return "uploadSetColor"
}
