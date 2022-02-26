package states

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type LookSetSeasonState struct{
	BaseState
}

func NewLookSetSeasonState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &LookSetSeasonState{BaseState: NewBase(bot, db, s3)}
}

func (s LookSetSeasonState) Do(message *tb.Message) string {
	s.db.SetUserSeason(message.Sender.ID, utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0])))

	return NewLookGenerateState(s.bot, s.db, s.s3).Do(message)
}

func (s LookSetSeasonState) GetName() string {
	return "lookSetSeason"
}
