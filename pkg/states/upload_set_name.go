package states

import (
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
)

type UploadSetNameState struct{
	BaseState
}

func NewUploadSetNameState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetNameState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetNameState) Do(message *tb.Message) string {
	s.db.SetName(message.Sender.ID, s.db.GetUser(message.Sender.ID).LastFileID, message.Text)
	s.bot.Send(message.Sender, constants.SendMeType, constants.TypeButtons())
	return UploadSetTypeState{}.GetName()
}

func (s UploadSetNameState) GetName() string {
	return "uploadSetName"
}
