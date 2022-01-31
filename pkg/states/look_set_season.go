package states

import (
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
)

type LookSetSeasonState struct{}

func (s LookSetSeasonState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	db.SetUserSeason(message.Sender.ID, strings.ToLower(strings.Split(message.Text, " ")[0]))

	return LookGenerateState{}.Do(bot, db, s3, message)
}

func (s LookSetSeasonState) GetName() string {
	return "lookSetSeason"
}
