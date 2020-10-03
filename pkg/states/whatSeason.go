package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"math/rand"
	"strings"
)

type WhatSeasonState struct{}

func (s WhatSeasonState) Do(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) string {
	topColor := db.GetTopColor(message.From.ID)
	bottomColor := db.GetBottomColor(message.From.ID)
	season := strings.ToLower(strings.Split(message.Text, " ")[0])
	db.SetUserSeason(message.From.ID, season)

	tops := GetByParams(db, message.From.ID, topColor, strings.ToLower(constants.Top), season)
	bottoms := GetByParams(db, message.From.ID, bottomColor, strings.ToLower(constants.Bottom), season)
	combos := GetByParams(db, message.From.ID, "any", strings.ToLower(constants.Combo), season)
	shoes := GetByParams(db, message.From.ID, "any", strings.ToLower(constants.Shoes), season)
	outer := GetByParams(db, message.From.ID, "any", strings.ToLower(constants.Outer), season)

	sepOrComb := rand.Intn(2)
	if sepOrComb == 0 && (len(tops) == 0 || len(bottoms) == 0) {
		sepOrComb = 1
	} else if sepOrComb == 1 && len(combos) == 0 {
		sepOrComb = 0
	}
	if sepOrComb == 0 {
		if len(tops) == 0 || len(bottoms) == 0 {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.NeedSomethingClean))
			return MainState{}.GetName()
		}
		SendSomething(bot, message, tops, false, strings.ToLower(constants.Top))
		SendSomething(bot, message, bottoms, false, strings.ToLower(constants.Bottom))
	} else {
		if len(combos) == 0 {
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.NeedSomethingClean))
			return MainState{}.GetName()
		}
		SendSomething(bot, message, combos, false, strings.ToLower(constants.Combo))
	}

	SendSomething(bot, message, shoes, false, strings.ToLower(constants.Shoes))
	SendSomething(bot, message, outer, false, strings.ToLower(constants.Outer))

	return MainState{}.GetName()
}

func (s WhatSeasonState) GetName() string {
	return "whatSeason"
}
