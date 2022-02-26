package states

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
)

type UploadSetSeasonState struct{
	BaseState
}

func NewUploadSetSeasonState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetSeasonState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetSeasonState) Do(message *tb.Message) string {
	_, err := s.bot.Send(message.Sender, constants.SendMeColor, constants.ColorButtons(false))
	if err != nil {
		log.Println("Err when send color request: ", err.Error())
	}

	return UploadSetColorState{}.GetName()
}

func (s UploadSetSeasonState) GetName() string {
	return "uploadSetSeason"
}
