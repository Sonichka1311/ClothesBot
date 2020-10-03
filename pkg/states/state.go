package states

import (
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
)

type State interface {
	Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string
	GetName() string
}

var (
	States = map[string]State{
		MainState{}.GetName():            MainState{},
		TypeState{}.GetName():            TypeState{},
		SeasonState{}.GetName():          SeasonState{},
		PhotoState{}.GetName():           PhotoState{},
		NameState{}.GetName():            NameState{},
		ColorState{}.GetName():           ColorState{},
		CategoryState{}.GetName():        CategoryState{},
		WhatTopColorState{}.GetName():    WhatTopColorState{},
		WhatBottomColorState{}.GetName(): WhatBottomColorState{},
		WhatSeasonState{}.GetName():      WhatSeasonState{},
	}
)
