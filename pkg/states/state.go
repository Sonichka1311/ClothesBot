package states

import (
	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
)

type State interface {
	Do(
		bot *tb.Bot,
		db *db.Database,
		s3 *s3.S3,
		message *tb.Message,
	) string
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
