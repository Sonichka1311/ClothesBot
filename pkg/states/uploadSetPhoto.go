package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UploadSetPhotoState struct{}

func (s UploadSetPhotoState) Do(bot *tb.Bot, db *db.Database, message *tb.Message) string {
	recent := db.GetRecent(message.Sender.ID)
	db.AddThing(message.Sender.ID, recent)
	db.SetPhoto(message.Sender.ID, recent, message.Photo.FileID)

	bot.Send(message.Sender, constants.SendMeName)
	return UploadSetNameState{}.GetName()
}

func (s UploadSetPhotoState) GetName() string {
	return "uploadSetPhoto"
}
