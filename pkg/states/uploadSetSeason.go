package states

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type UploadSetSeasonState struct{}

func (s UploadSetSeasonState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	_, err := bot.Send(message.Sender, constants.SendMeColor, constants.ColorButtons(false))
	if err != nil {
		log.Println("Err when send color request: ", err.Error())
	}

	return UploadSetColorState{}.GetName()
}

func (s UploadSetSeasonState) GetName() string {
	return "uploadSetSeason"
}
