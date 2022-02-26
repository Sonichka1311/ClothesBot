package states

import (
	"math/rand"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type LookGenerateState struct{
	BaseState
}

func NewLookGenerateState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &LookGenerateState{BaseState: NewBase(bot, db, s3)}
}

func (s LookGenerateState) Do(message *tb.Message) string {
	uid := message.Sender.ID
	user := s.db.GetUser(uid)

	season := utils.ToEng(strings.ToLower(strings.Split(message.Text, " ")[0]))

	tops := s.db.GetByParams(uid, user.TopColor, utils.ToEng(constants.Top), season)
	bottoms := s.db.GetByParams(uid, user.BottomColor, utils.ToEng(constants.Bottom), season)
	combos := s.db.GetByParams(uid, "any", utils.ToEng(constants.Combo), season)
	shoes := s.db.GetByParams(uid, "any", utils.ToEng(constants.Shoes), season)
	outer := s.db.GetByParams(uid, "any", utils.ToEng(constants.Outer), season)

	photos := make([]tb.InputMedia, 0)
	texts := make([]string, 0)

	buttons := make([]*tb.InlineButton, 0)

	look := &constants.Look{}

	sepOrComb := rand.Intn(2)
	if sepOrComb == 0 && (len(tops) == 0 || len(bottoms) == 0) {
		sepOrComb = 1
	} else if sepOrComb == 1 && len(combos) == 0 {
		sepOrComb = 0
	}

	if sepOrComb == 0 {
		if len(tops) == 0 || len(bottoms) == 0 {
			s.bot.Send(message.Sender, constants.NeedSomethingClean)
			return MainState{}.GetName()
		}

		top := GenerateSomething(tops)
		s.appendSomething(top, &texts, &photos, &buttons)
		look.BaseTop = top

		bottom := GenerateSomething(bottoms)
		s.appendSomething(bottom, &texts, &photos, &buttons)
		look.Bottom = bottom
	} else {
		if len(combos) == 0 {
			s.bot.Send(message.Sender, constants.NeedSomethingClean)
			return MainState{}.GetName()
		}

		combo := GenerateSomething(combos)
		s.appendSomething(combo, &texts, &photos, &buttons)
		look.Combination = combo
	}

	shoe := GenerateSomething(shoes)
	s.appendSomething(shoe, &texts, &photos, &buttons)
	look.Shoes = shoe

	outerwear := GenerateSomething(outer)
	s.appendSomething(outerwear, &texts, &photos, &buttons)
	look.Outerwear = outerwear

	s.sendAlbumAndCaption(
		message.Sender,
		look,
		&texts, &photos, &buttons,
		sepOrComb == 0,
		(sepOrComb == 0 && len(combos) != 0) || (sepOrComb == 1 && len(tops) != 0 && len(bottoms) != 0),
	)

	return MainState{}.GetName()
}

func (s LookGenerateState) GetName() string {
	return "lookGenerate"
}
