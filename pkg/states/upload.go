package states

import (
	"strconv"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type CreateState struct{}

func (s CreateState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	recent, _ := strconv.Atoi(message.Text)
	bot.Send(message.Sender, constants.Added(db.GetThing(message.Sender.ID, recent).Name, recent))
	return MainState{}.GetName()
}

func (s CreateState) GetName() string {
	return "waitColor"
}

