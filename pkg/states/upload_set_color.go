package states

import (
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type UploadSetColorState struct{
	BaseState
}

func NewUploadSetColorState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetColorState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetColorState) Do(message *tb.Message) string {
	recent := s.db.GetUser(message.Sender.ID).LastFileID
	s.db.SetColor(message.Sender.ID, recent, utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0])))
	message.Text = strconv.Itoa(recent)
	s.db.SetRecent(message.Sender.ID, recent+1)
	return NewCreateState(s.bot, s.db, s.s3).Do(message)
}

func (s UploadSetColorState) GetName() string {
	return "uploadSetColor"
}
