package states

import (
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/db"
	"bot/pkg/s3"
)

type State interface {
	Do(message *tb.Message) string
	GetName() string
}

type StateFabric struct {
	bot *tb.Bot
	db  *db.Database
	s3  *s3.S3
}

func NewStateFabric(bot *tb.Bot, db *db.Database, s3 *s3.S3) *StateFabric {
	return &StateFabric{
		bot: bot,
		db:  db,
		s3:  s3,
	}
}

func (f StateFabric) NewState(state string) State {
	return States[state](f.bot, f.db, f.s3)
}

type NewStateFunc func(*tb.Bot, *db.Database, *s3.S3) State

var (
	States = map[string]NewStateFunc{
		BaseState{}.GetName(): NewBaseState,
		MainState{}.GetName(): NewMainState,

		UploadSetPhotoState{}.GetName():     NewUploadSetPhotoState,
		UploadSetNameState{}.GetName():      NewUploadSetNameState,
		UploadSetTypeState{}.GetName():      NewUploadSetTypeState,
		UploadSetComboState{}.GetName():     NewUploadSetComboState,
		UploadSetComboTypeState{}.GetName(): NewUploadSetComboTypeState,
		UploadSetSeasonState{}.GetName():    NewUploadSetSeasonState,
		UploadSetColorState{}.GetName():     NewUploadSetColorState,

		LookGenerateState{}.GetName():       NewLookGenerateState,
		LookSetTopColorState{}.GetName():    NewLookSetTopColorState,
		LookSetBottomColorState{}.GetName(): NewLookSetBottomColorState,
		LookSetSeasonState{}.GetName():      NewLookSetSeasonState,
	}
)
