package states

import (
	"strings"
	"sync"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type UploadSetTypeState struct{
	BaseState
}

func NewUploadSetTypeState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetTypeState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetTypeState) Do(message *tb.Message) string {
	thingType := utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0]))
	s.db.SetType(message.Sender.ID, s.db.GetUser(message.Sender.ID).LastFileID, thingType)

	constants.Mutex.Lock()
	constants.MutexMap[message.Sender.ID] = &sync.Mutex{}
	constants.Mutex.Unlock()

	//if thingType != "top" {
		s.bot.Send(message.Sender, constants.SendMeSeason, constants.SeasonButtons(true))
		return UploadSetSeasonState{}.GetName()
	//}
	//
	//s.bot.Send(message.Sender, constants.SendMeSeason, constants.ComboButtons(true))
	//return UploadSetComboState{}.GetName()
}

func (s UploadSetTypeState) GetName() string {
	return "uploadSetType"
}
