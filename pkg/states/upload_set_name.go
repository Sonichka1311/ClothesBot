package states

import (
	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type UploadSetNameState struct{}

func (s UploadSetNameState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	db.SetName(message.Sender.ID, db.GetUser(message.Sender.ID).LastFileID, message.Text)
	bot.Send(message.Sender, constants.SendMeType, constants.TypeButtons())
	return UploadSetTypeState{}.GetName()
}

func (s UploadSetNameState) GetName() string {
	return "uploadSetName"
}
