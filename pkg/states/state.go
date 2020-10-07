package states

import (
	"bot/pkg/db"
	tb "gopkg.in/tucnak/telebot.v2"
)

type State interface {
	Do(bot *tb.Bot, db *db.Database, message *tb.Message) string
	GetName() string
}

var (
	States = map[string]State{
		MainState{}.GetName():               MainState{},
		UploadSetTypeState{}.GetName():      UploadSetTypeState{},
		UploadSetSeasonState{}.GetName():    UploadSetSeasonState{},
		UploadSetPhotoState{}.GetName():     UploadSetPhotoState{},
		UploadSetNameState{}.GetName():      UploadSetNameState{},
		UploadSetColorState{}.GetName():     UploadSetColorState{},
		LookSetTopColorState{}.GetName():    LookSetTopColorState{},
		LookSetBottomColorState{}.GetName(): LookSetBottomColorState{},
		LookSetSeasonState{}.GetName():      LookSetSeasonState{},
	}
)
