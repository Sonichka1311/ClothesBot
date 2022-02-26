package states

import (
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
	"bot/pkg/s3"
)

type UploadSetComboState struct{
	BaseState
}

func NewUploadSetComboState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetComboTypeState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetComboState) Do(message *tb.Message) string {
	//s.db.SetType(message.Sender.ID, s.db.GetUser(message.Sender.ID).LastFileID, strings.ToLower(strings.Split(message.Text, " ")[0]))
	//
	//constants.Mutex.Lock()
	//constants.MutexMap[message.Sender.ID] = &sync.Mutex{}
	//constants.Mutex.Unlock()

	//s.bot.Send(message.Sender, constants.SendMeSeason, constants.SeasonButtons(true))
	return UploadSetComboTypeState{}.GetName()
}

func (s UploadSetComboState) GetName() string {
	return "uploadSetComboState"
}

