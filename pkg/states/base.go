package states

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/utils"
)

type BaseState struct {
	bot *tb.Bot
	db  *db.Database
	s3  *s3.S3
}

func NewBase(bot *tb.Bot, db *db.Database, s3 *s3.S3) BaseState {
	return BaseState{
		bot: bot,
		db:  db,
		s3:  s3,
	}
}

func NewBaseState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &BaseState{
		bot: bot,
		db:  db,
		s3:  s3,
	}
}

func (s BaseState) Do(message *tb.Message) string {
	switch message.Text {
	case "/start":
		s.db.CreateUser(message.Sender.ID)
		s.hello(message.Sender)
		return MainState{}.GetName()
	case "/end":
		user := s.db.GetUser(message.Sender.ID)
		if strings.HasPrefix(user.State, "upload") {
			s.deleteThing(message, strconv.Itoa(user.LastFileID), user.State != UploadSetPhotoState{}.GetName())
			s.bot.Send(message.Sender, constants.UploadCancelled)
		}
		s.remind(message.Sender)
		return MainState{}.GetName()
	case "/help":
		s.help(message.Sender)
		return MainState{}.GetName()
	}
	return s.GetName()
}

func (s BaseState) GetName() string {
	return "base"
}

func (s BaseState) remind(sender *tb.User) {
	s.bot.Send(sender, constants.Remind)
}

func (s BaseState) hello(sender *tb.User) {
	s.bot.Send(sender, constants.Hello)
}

func (s BaseState) help(sender *tb.User) {
	s.bot.Send(sender, constants.Help)
}

func (s BaseState) smthWrong(sender *tb.User) {
	s.bot.Send(sender, constants.CommandNotFound)
}

func (s BaseState) sendBigMessage(sender *tb.User, text string) {
	texts := SplitBigMsg(text)
	for _, message := range texts {
		s.bot.Send(sender, message)
	}
}

func (s BaseState) appendSomething(
	thing *constants.Thing,
	texts *[]string,
	photos *[]tb.InputMedia,
	buttons *[]*tb.InlineButton,
) {
	if thing.Photo == "" {
		*texts = append(*texts, constants.NoCleanThing[thing.Type])
		return
	}

	photo := s.createMediaByThing(thing, thing.ShortCaption)
	*texts = append(*texts, photo.Caption)

	if photo.File.FileID != "" || photo.File.FileReader != nil {
		photo.Caption = strings.Split(photo.Caption, "\n")[0]
		*photos = append(*photos, photo)
		*buttons = append(*buttons, constants.ChangeButton(thing.Type))
	}
}

func (s BaseState) ChangeThing(
	data string,
	message *tb.Message,
	types string,
	change bool,
) {
	if change {
		user := s.db.GetUser(message.Sender.ID)

		colors := map[string]string{
			strings.ToLower(constants.Top):    user.TopColor,
			strings.ToLower(constants.Bottom): user.BottomColor,
		}

		color := "any"
		if c, ok := colors[types]; ok {
			color = c
		}

		photoMsgId, _ := strconv.Atoi(strings.Split(data, "_")[1])
		s.changeSomething(message, s.db.GetByParams(message.Sender.ID, color, types, user.Season), types, photoMsgId)
	} else {
		s.changeType(message, types, strings.Split(data, "_")[2:])
	}
}

func (s BaseState) changeSomething(
	message *tb.Message,
	things []*constants.Thing,
	thingType string,
	photoMessageID int,
) {
	if len(things) == 0 {
		s.bot.Send(message.Sender, constants.NoCleanThing[thingType])
		return
	}

	thing := things[rand.Intn(len(things))]
	texts := s.rebuildLookText(message.Text, message.Sender.ID)
	textIdx := GetTextIndex(thingType, len(texts))

	if len(things) > 1 {
		ReGenerateThing(GetThingIdByCaption(texts[textIdx]), &thing, things)
	} else {
		texts[textIdx] = thing.ShortCaption() + constants.JustOneThing
		buttonIdx := GetButtonIndex(textIdx, texts)
		keyboard := GetKeyboardWithoutButtonWithIdx(message.ReplyMarkup.InlineKeyboard, buttonIdx)
		msg, err := s.bot.Edit(message, strings.Join(texts, "\n"), keyboard)
		if err != nil {
			log.Println("Err edit capt just one thing: ", err.Error())
			return
		}

		go s.hideButtons(msg)
		return
	}

	photo := s.createMediaByThing(
		thing,
		func() string {
			return strings.Split(thing.ShortCaption(), "\n")[0]
		},
	)
	_, err := s.bot.EditMedia(&tb.Message{ID: photoMessageID, Chat: message.Chat}, photo)
	if err != nil {
		log.Println("Err edit media: " + err.Error())
		return
	}

	keyboard := tb.ReplyMarkup{InlineKeyboard: message.ReplyMarkup.InlineKeyboard}

	msg, err := s.bot.Edit(
		message,
		ReplaceThingText(texts, thing, textIdx),
		&keyboard,
	)
	if err != nil {
		log.Println("Err edit text: " + err.Error())
		return
	}

	go s.hideButtons(msg)
}

func (s BaseState) changeType(
	message *tb.Message,
	thingType string,
	thingsIDs []string,
) {
	texts := make([]string, 0)
	photos := make([]tb.InputMedia, 0)
	buttons := make([]*tb.InlineButton, 0)

	uid := message.Sender.ID
	user := s.db.GetUser(uid)
	season := user.Season

	look := &constants.Look{}

	if thingType == "sep" {
		topColor := user.TopColor
		bottomColor := user.BottomColor
		tops := s.db.GetByParams(uid, topColor, strings.ToLower(constants.Top), season)
		bottoms := s.db.GetByParams(uid, bottomColor, strings.ToLower(constants.Bottom), season)

		top := GenerateSomething(tops)
		s.appendSomething(top, &texts, &photos, &buttons)
		look.BaseTop = top

		bottom := GenerateSomething(bottoms)
		s.appendSomething(bottom, &texts, &photos, &buttons)
		look.Bottom = bottom
	} else if thingType == "comb" {
		combos := s.db.GetByParams(uid, "any", strings.ToLower(constants.Combo), season)

		combo := GenerateSomething(combos)
		s.appendSomething(combo, &texts, &photos, &buttons)
		look.Combination = combo
	}

	for idx, thingID := range thingsIDs {
		id, _ := strconv.Atoi(thingID)
		if id != -1 {
			thing := s.db.GetThing(uid, id)
			s.appendSomething(thing, &texts, &photos, &buttons)
			if idx == 0 {
				look.Shoes = thing
			} else {
				look.Outerwear = thing
			}
		} else {
			if idx == 0 {
				texts = append(texts, constants.NoCleanThing[strings.ToLower(constants.Shoes)])
			} else {
				texts = append(texts, constants.NoCleanThing[strings.ToLower(constants.Outer)])
			}
		}
	}

	s.sendAlbumAndCaption(message.Sender, look, &texts, &photos, &buttons, thingType == "sep", true)
}

func (s BaseState) rebuildLookText(text string, userId int) []string {
	texts := strings.Split(text, "\n\n")

	for idx, text := range texts {
		stringsInText := strings.Split(text, "\n")

		if len(stringsInText) > 1 {
			id := GetThingIdByCaption(text)

			if id != -1 {
				texts[idx] = s.db.GetThing(userId, id).ShortCaption()

				if len(stringsInText) > 2 {
					texts[idx] += strings.Join(append(stringsInText[2:]), "\n") + "\n"
				}
			}
		}
	}

	return texts
}

func (s BaseState) sendAlbumAndCaption(
	sender *tb.User,
	look *constants.Look,
	texts *[]string,
	photos *[]tb.InputMedia,
	buttons *[]*tb.InlineButton,
	isSeparate bool,
	haveDifferentType bool,
) {
	img := imaging.New(2500, 3000, image.Transparent)

	end := 0
	main := imaging.New(2500, 3000, image.Transparent)
	if isSeparate {
		top, _ := png.Decode((*photos)[0].MediaFile().FileReader)
		bottom, _ := png.Decode((*photos)[1].MediaFile().FileReader)

		topModified := ResizeImage(CropImage(top), 600)
		bottomModified := ResizeImage(CropImage(bottom), 600)

		main = imaging.Overlay(main, bottomModified, image.Point{X: 800, Y: topModified.Bounds().Dy() * 9 / 10}, 1)
		main = imaging.Overlay(main, topModified, image.Point{X: 800}, 1)

		end = 2
	} else {
		combo, _ := png.Decode((*photos)[0].MediaFile().FileReader)
		comboModified := ResizeImage(CropImage(combo), 600)

		main = imaging.Overlay(main, comboModified, image.Point{X: 800}, 1)

		end = 1
	}

	if look.Shoes != nil {
		shoes, _ := png.Decode((*photos)[end].MediaFile().FileReader)
		shoesModified := imaging.Resize(CropImage(shoes), 600, 0, imaging.NearestNeighbor)

		main = imaging.Overlay(main, shoesModified, image.Point{X: 1500, Y: 1700}, 1)

		end++
	}

	if look.Outerwear != nil {
		outer, _ := png.Decode((*photos)[end].MediaFile().FileReader)
		outerModified := ResizeImage(CropImage(outer), 800)

		img = imaging.Overlay(img, outerModified, image.Point{}, 1)
	}

	img = imaging.Overlay(img, main, image.Point{}, 1)


	var out bytes.Buffer
	png.Encode(&out, img)

	media := &tb.Photo{
		File:      tb.File{FileReader: &out},
		ParseMode: constants.ParseMode,
		Caption:   strings.Join(*texts, "\n"),
	}
	//s.bot.Send(sender, media)

	//messages, _ := s.bot.SendAlbum(sender, *photos)
	//for idx, message := range messages {
	//	(*buttons)[idx].Data = strings.Split((*buttons)[idx].Data, " ")[0] + "_" + strconv.Itoa(message.ID)
	//}

	if haveDifferentType {
		startIndex := 1
		if isSeparate {
			*buttons = append(*buttons, constants.ChangeTypeButton("sep", "comb"))
			startIndex = 2
		} else {
			*buttons = append(*buttons, constants.ChangeTypeButton("comb", "sep"))
		}

		for idx := startIndex; idx < len(*texts); idx++ {
			(*buttons)[len(*buttons)-1].Data += fmt.Sprintf("_%d", GetThingIdByCaption((*texts)[idx]))
		}
	}

	keyboard := constants.NewKeyboard(*buttons...)
	msg, _ := s.bot.Send(sender, media, keyboard)
	go s.hideButtons(msg)
}

var ToHide = map[int]int{} //messageID to number of changes

func (s BaseState) hideButtons(message *tb.Message) {
	ToHide[message.ID]++

	time.Sleep(time.Minute * 5)
	ToHide[message.ID]--

	if ToHide[message.ID] == 0 {
		s.bot.EditCaption(message, message.Caption+constants.TimeIsUp)
		delete(ToHide, message.ID)
	} else if ToHide[message.ID] < 0 {
		delete(ToHide, message.ID)
	}
}

func (s BaseState) MultiCallback(
	message *tb.Callback,
	check, uncheck func(id int, recent int, s string),
) {
	constants.MutexMap[message.Sender.ID].Lock()
	defer constants.MutexMap[message.Sender.ID].Unlock()

	keyboard := message.Message.ReplyMarkup
	buttons := make([][]tb.InlineButton, 0)

	user := s.db.GetUser(message.Sender.ID)

	for _, button := range keyboard.InlineKeyboard {
		newButton := []tb.InlineButton{*constants.NewButton(message.Data, button[0].Text)}
		text := utils.ToEng(strings.ToLower(strings.Split(button[0].Text, " ")[0]))

		if message.Data != constants.Done {
			switch {
			case strings.HasPrefix(message.Data, button[0].Text):
				buttons = append(buttons, newButton)
				check(message.Sender.ID, user.LastFileID, text)
			case strings.HasPrefix(button[0].Text, message.Data):
				buttons = append(buttons, newButton)
				uncheck(message.Sender.ID, user.LastFileID, text)
			default:
				buttons = append(buttons, []tb.InlineButton{button[0]})
			}
		} else {
			buttons = append(buttons, []tb.InlineButton{button[0]})
		}
	}

	newKeyboard := tb.ReplyMarkup{InlineKeyboard: buttons}
	s.bot.EditReplyMarkup(message.Message, &newKeyboard)
}

func (s BaseState) createMediaByThing(thing *constants.Thing, caption func() string) *tb.Photo {
	res, err := s.s3.GetObject(thing.Photo)
	if err != nil {
		log.Println(err)
	}

	return &tb.Photo{
		File:      tb.File{FileReader: res},
		Caption:   caption(),
		ParseMode: constants.ParseMode,
	}
}

func (s BaseState) deleteThing(message *tb.Message, thingIDStr string, all bool) {
	thingID, _ := strconv.Atoi(thingIDStr)
	thing := s.db.GetThing(message.Sender.ID, thingID)

	if all {
		s.db.DeleteThing(message.Sender.ID, thingID)
		err := s.s3.DeleteObject(thing.Photo)
		if err != nil {
			log.Println(err)
		}
	}
}
