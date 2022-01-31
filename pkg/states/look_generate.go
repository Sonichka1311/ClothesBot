package states

import (
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type LookGenerateState struct{}

func (s LookGenerateState) Do(bot *tb.Bot, db *db.Database, s3 *s3.S3, message *tb.Message) string {
	uid := message.Sender.ID
	user := db.GetUser(uid)

	season := strings.ToLower(strings.Split(message.Text, " ")[0])

	tops := db.GetByParams(uid, user.TopColor, strings.ToLower(constants.Top), season)
	bottoms := db.GetByParams(uid, user.BottomColor, strings.ToLower(constants.Bottom), season)
	combos := db.GetByParams(uid, "any", strings.ToLower(constants.Combo), season)
	shoes := db.GetByParams(uid, "any", strings.ToLower(constants.Shoes), season)
	outer := db.GetByParams(uid, "any", strings.ToLower(constants.Outer), season)

	photos := make([]tb.InputMedia, 0)
	texts := make([]string, 0)

	buttons := make([]*tb.InlineButton, 0)

	sepOrComb := rand.Intn(2)
	if sepOrComb == 0 && (len(tops) == 0 || len(bottoms) == 0) {
		sepOrComb = 1
	} else if sepOrComb == 1 && len(combos) == 0 {
		sepOrComb = 0
	}

	if sepOrComb == 0 {
		if len(tops) == 0 || len(bottoms) == 0 {
			bot.Send(message.Sender, constants.NeedSomethingClean)
			return MainState{}.GetName()
		}

		AppendSomething(s3, GenerateSomething(tops, strings.ToLower(constants.Top)), &texts, &photos, &buttons)
		AppendSomething(s3, GenerateSomething(bottoms, strings.ToLower(constants.Bottom)), &texts, &photos, &buttons)
	} else {
		if len(combos) == 0 {
			bot.Send(message.Sender, constants.NeedSomethingClean)
			return MainState{}.GetName()
		}

		AppendSomething(s3, GenerateSomething(combos, strings.ToLower(constants.Combo)), &texts, &photos, &buttons)
	}

	AppendSomething(s3, GenerateSomething(shoes, strings.ToLower(constants.Shoes)), &texts, &photos, &buttons)
	AppendSomething(s3, GenerateSomething(outer, strings.ToLower(constants.Outer)), &texts, &photos, &buttons)

	SendAlbumAndCaption(bot, message.Sender,
		&texts, &photos, &buttons,
		sepOrComb == 0,
		(sepOrComb == 0 && len(combos) != 0) || (sepOrComb == 1 && len(tops) != 0 && len(bottoms) != 0),
	)

	return MainState{}.GetName()
}

func (s LookGenerateState) GetName() string {
	return "lookGenerate"
}
