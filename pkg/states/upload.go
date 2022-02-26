package states

import (
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
)

type CreateState struct{
	BaseState
}

func NewCreateState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &CreateState{BaseState: NewBase(bot, db, s3)}
}

func (s CreateState) Do(message *tb.Message) string {
	recent, _ := strconv.Atoi(message.Text)
	s.bot.Send(message.Sender, constants.Added(s.db.GetThing(message.Sender.ID, recent).Name, recent))
	return MainState{}.GetName()
}

func (s CreateState) GetName() string {
	return "waitColor"
}

