package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Remind(bot *tb.Bot, sender *tb.User) {
	go bot.Send(sender, constants.Remind)
}

func Hello(bot *tb.Bot, sender *tb.User) {
	go bot.Send(sender, constants.Hello)
}

func Help(bot *tb.Bot, sender *tb.User) {
	go bot.Send(sender, constants.Help)
}

func SmthWrong(bot *tb.Bot, sender *tb.User) {
	go bot.Send(sender, constants.CommandNotFound)
}

func GetThing(bot *tb.Bot, db *db.Database, message *tb.Message) {
	id, _ := strconv.Atoi(message.Text[7:])
	thing := db.GetThing(message.Sender.ID, id)
	photo := CreateMediaByThing(thing, constants.ThingText)
	bot.Send(message.Sender, photo)
}

func Wardrobe(bot *tb.Bot, db *db.Database, message *tb.Message) {
	rows := db.GetAll(message.Sender.ID)
	texts := make([]string, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Purity)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe from database: %s\n", dbError.Error())
		}
		texts = append(texts, constants.ThingsText(&thing))
	}
	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		bot.Send(message.Sender, constants.EmptyArray["wardrobe"])
	}
	SendBigMsg(bot, message.Sender, text)
}

func ChangePurity(bot *tb.Bot, db *db.Database, message *tb.Message, toClean bool) {
	id, _ := strconv.Atoi(message.Text[7:])
	if toClean {
		db.MakeClean(message.Sender.ID, id)
	} else {
		db.MakeDirty(message.Sender.ID, id)
	}
	bot.Send(message.Sender, constants.MarkThingTo(db.GetName(message.Sender.ID, id), toClean)+message.Text[7:])
}

func Dirty(bot *tb.Bot, db *db.Database, message *tb.Message) {
	rows := db.GetDirty(message.Sender.ID)
	texts := make([]string, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name)
		thing.Purity = "dirty"
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe from database: %s\n", dbError.Error())
		}
		texts = append(texts, constants.ThingsText(&thing))
	}
	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		bot.Send(message.Sender, constants.EmptyArray["dirty"])
	}
	SendBigMsg(bot, message.Sender, text)
}

func GetByType(bot *tb.Bot, db *db.Database, message *tb.Message) {
	rows := db.GetByType(message.Sender.ID, message.Text[6:])
	texts := make([]string, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Purity, &thing.Photo)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe by type from database: %s\n", dbError.Error())
		}
		texts = append(texts, constants.ThingsText(&thing))
	}
	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		bot.Send(message.Sender, constants.EmptyArray["by_type"])
	}
	SendBigMsg(bot, message.Sender, text)
}

func SendBigMsg(bot *tb.Bot, sender *tb.User, text string) {
	texts := constants.SplitBigMsg(text)
	for _, message := range texts {
		bot.Send(sender, message)
	}
}

func GetByParams(db *db.Database, id int, color string, types string, season string) []*constants.Thing {
	ans := make([]*constants.Thing, 0)
	rows := db.GetByParams(id, color, types, season)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Photo)
		thing.Type = types
		if dbError != nil {
			log.Printf("Error while selecting %ss from database: %s\n", types, dbError.Error())
		}
		ans = append(ans, &thing)
	}
	return ans
}

func AppendSomething(
	thing *constants.Thing,
	texts *[]string,
	photos *[]tb.InputMedia,
	buttons *[]*tb.InlineButton,
) {
	if thing.Photo == "" {
		*texts = append(*texts, constants.NoCleanThing[thing.Type])
		return
	}
	photo := CreateMediaByThing(thing, constants.Caption)
	*texts = append(*texts, photo.Caption)
	if photo.File.FileID != "" {
		photo.Caption = strings.Split(photo.Caption, "\n")[0]
		*photos = append(*photos, photo)
		*buttons = append(*buttons, constants.ChangeButton(thing.Type))
	}
}

func GenerateSomething(things []*constants.Thing, thingType string) *constants.Thing {
	if len(things) == 0 {
		return &constants.Thing{Type: thingType}
	}
	return things[rand.Intn(len(things))]
}

func GetTextIndex(thingType string, numberElements int) int {
	idx := 0
	switch thingType {
	case "bottom":
		idx = 1
	case "shoes":
		idx = 2
		if numberElements == 3 {
			idx = 1
		}
	case "outer":
		idx = 3
		if numberElements == 3 {
			idx = 2
		}
	}
	return idx
}

func GetThingIdByCaption(text string) int {
	stringsInText := strings.Split(text, "\n")
	if len(stringsInText) > 1 {
		num, err := strconv.Atoi(strings.Split(stringsInText[1], "_")[1])
		if err != nil {
			log.Println("Can't get id by caption: ", err.Error())
			return -1
		}
		return num
	}
	return -1
}

func RebuildLookText(db *db.Database, text string, userId int) []string {
	texts := strings.Split(text, "\n\n")
	for idx, text := range texts {
		stringsInText := strings.Split(text, "\n")
		if len(stringsInText) > 1 {
			id := GetThingIdByCaption(text)
			if id != -1 {
				texts[idx] = constants.Caption(db.GetThing(userId, id))
				if len(stringsInText) > 2 {
					texts[idx] += strings.Join(append(stringsInText[2:]), "\n") + "\n"
				}
			}
		}
	}
	return texts
}

func ReGenerateThing(oldThingId int, thing **constants.Thing, things []*constants.Thing) {
	for oldThingId == (*thing).Id {
		log.Println("The same thing occurred.")
		*thing = things[rand.Intn(len(things))]
	}
}

func GetKeyboardWithoutButtonWithIdx(keyboard [][]tb.InlineButton, idx int) *tb.ReplyMarkup {
	return &tb.ReplyMarkup{
		InlineKeyboard: append(keyboard[:idx], keyboard[idx+1:]...),
	}
}

func GetButtonIndex(textIdx int, texts []string) int {
	buttonIdx := textIdx
	for idx := 0; idx < textIdx; idx++ {
		if len(strings.Split(texts[idx], "\n")) > 3 {
			buttonIdx--
		}
	}
	return buttonIdx
}

func ReplaceThingText(texts []string, thing *constants.Thing, textIdx int) string {
	return strings.Join(
		append(
			append(texts[:textIdx], constants.Caption(thing)),
			texts[textIdx+1:]...,
		),
		"\n",
	)
}

func SendAlbumAndCaption(
	bot *tb.Bot,
	sender *tb.User,
	texts *[]string,
	photos *[]tb.InputMedia,
	buttons *[]*tb.InlineButton,
	isSeparate bool,
	haveDifferentType bool,
) {
	messages, _ := bot.SendAlbum(sender, *photos)
	for idx, message := range messages {
		(*buttons)[idx].Data = strings.Split((*buttons)[idx].Data, " ")[0] + "_" + strconv.Itoa(message.ID)
	}
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
	msg, _ := bot.Send(sender, strings.Join(*texts, "\n"), keyboard)
	go HideButtons(bot, msg)
}

func ChangeType(bot *tb.Bot, db *db.Database, message *tb.Message, thingType string, thingsIDs []string) {
	texts := make([]string, 0)
	photos := make([]tb.InputMedia, 0)
	buttons := make([]*tb.InlineButton, 0)

	uid := message.Sender.ID
	season := db.GetSeason(uid)
	if thingType == "sep" {
		topColor := db.GetTopColor(uid)
		bottomColor := db.GetBottomColor(uid)
		tops := GetByParams(db, uid, topColor, strings.ToLower(constants.Top), season)
		bottoms := GetByParams(db, uid, bottomColor, strings.ToLower(constants.Bottom), season)
		AppendSomething(GenerateSomething(tops, strings.ToLower(constants.Top)), &texts, &photos, &buttons)
		AppendSomething(GenerateSomething(bottoms, strings.ToLower(constants.Bottom)), &texts, &photos, &buttons)
	} else if thingType == "comb" {
		combos := GetByParams(db, uid, "any", strings.ToLower(constants.Combo), season)
		AppendSomething(GenerateSomething(combos, strings.ToLower(constants.Combo)), &texts, &photos, &buttons)
	}

	for idx, thingID := range thingsIDs {
		id, _ := strconv.Atoi(thingID)
		if id != -1 {
			thing := db.GetThing(uid, id)
			AppendSomething(thing, &texts, &photos, &buttons)
		} else {
			if idx == 0 {
				texts = append(texts, constants.NoCleanThing[strings.ToLower(constants.Shoes)])
			} else {
				texts = append(texts, constants.NoCleanThing[strings.ToLower(constants.Outer)])
			}
		}
	}

	SendAlbumAndCaption(bot, message.Sender, &texts, &photos, &buttons, thingType == "sep", true)
}

func ChangeSomething(bot *tb.Bot, db *db.Database, message *tb.Message, things []*constants.Thing, thingType string, photoMessageID int) {
	if len(things) == 0 {
		bot.Send(message.Sender, constants.NoCleanThing[thingType])
		return
	}
	thing := things[rand.Intn(len(things))]
	texts := RebuildLookText(db, message.Text, message.Sender.ID)
	textIdx := GetTextIndex(thingType, len(texts))
	if len(things) > 1 {
		ReGenerateThing(GetThingIdByCaption(texts[textIdx]), &thing, things)
	} else {
		texts[textIdx] = constants.Caption(thing)+constants.JustOneThing
		buttonIdx := GetButtonIndex(textIdx, texts)
		keyboard := GetKeyboardWithoutButtonWithIdx(message.ReplyMarkup.InlineKeyboard, buttonIdx)
		msg, err := bot.Edit(message, strings.Join(texts, "\n"), keyboard)
		if err != nil {
			log.Println("Err edit capt just one thing: ", err.Error())
			return
		}
		go HideButtons(bot, msg)
		return
	}

	photo := CreateMediaByThing(
		thing,
		func (thing *constants.Thing) string {
			return strings.Split(constants.Caption(thing), "\n")[0]
		},
	)
	_, err := bot.EditMedia(&tb.Message{ID: photoMessageID, Chat: message.Chat}, photo)
	if err != nil {
		log.Println("Err edit media: " + err.Error())
		return
	}

	keyboard := tb.ReplyMarkup{InlineKeyboard: message.ReplyMarkup.InlineKeyboard}

	msg, err := bot.Edit(
		message,
		ReplaceThingText(texts, thing, textIdx),
		&keyboard,
	)
	if err != nil {
		log.Println("Err edit text: " + err.Error())
		return
	}

	go HideButtons(bot, msg)
}

var ToHide = map[int]int{} //messageID to number of changes

func HideButtons(bot *tb.Bot, message *tb.Message) {
	if _, ok := ToHide[message.ID]; ok {
		ToHide[message.ID] += 1
	} else {
		ToHide[message.ID] = 1
	}
	time.Sleep(time.Minute * 5)
	ToHide[message.ID] -= 1
	if ToHide[message.ID] == 0 {
		bot.EditCaption(message, message.Caption+constants.TimeIsUp)
	}
}

func ChangeThing(bot *tb.Bot, db *db.Database, data string, message *tb.Message, types string, change bool) {
	if change {
		topColor := db.GetTopColor(message.Sender.ID)
		bottomColor := db.GetBottomColor(message.Sender.ID)
		colors := map[string]string{
			strings.ToLower(constants.Top):    topColor,
			strings.ToLower(constants.Bottom): bottomColor,
		}
		season := db.GetSeason(message.Sender.ID)
		color := "any"
		if c, ok := colors[types]; ok {
			color = c
		}
		photoMsgId, _ := strconv.Atoi(strings.Split(data, "_")[1])
		ChangeSomething(bot, db, message, GetByParams(db, message.Sender.ID, color, types, season), types, photoMsgId)
	} else {
		ChangeType(bot, db, message, types, strings.Split(data, "_")[2:])
	}
}

func GetRandomThing(bot *tb.Bot, db *db.Database, message *tb.Message, types string) {
	rows := db.GetByType(message.Sender.ID, types)
	ans := make([]*constants.Thing, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Purity, &thing.Photo)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe by type from database: %s\n", dbError.Error())
		}
		thing.Type = types
		if thing.Purity == "clean" {
			ans = append(ans, &thing)
		}
	}
	if len(ans) == 0 {
		bot.Send(message.Sender, constants.EmptyArray["random"])
		return
	}
	thing := ans[rand.Intn(len(ans))]
	photo := CreateMediaByThing(thing, constants.Caption)
	bot.Send(message.Sender, photo)
	log.Println(photo.Caption)
}

func DeleteThing(bot *tb.Bot, db *db.Database, message *tb.Message, idStr string) {
	id, _ := strconv.Atoi(idStr)
	go db.DeleteThing(message.Sender.ID, id)
	bot.Send(message.Sender, constants.Deleted)
}

func MultiCallback(
	bot *tb.Bot,
	db *db.Database,
	message *tb.Callback,
	do, undo func(id int, recent int, s string),
) {
	constants.MutexMap[message.Sender.ID].Lock()
	defer constants.MutexMap[message.Sender.ID].Unlock()
	keyboard := message.Message.ReplyMarkup
	buttons := make([][]tb.InlineButton, 0)
	for _, button := range keyboard.InlineKeyboard {
		newButton := []tb.InlineButton{*constants.NewButton(message.Data, button[0].Text)}
		text := strings.ToLower(strings.Split(button[0].Text, " ")[0])
		if strings.HasPrefix(message.Data, button[0].Text) && message.Data != "Done" {
			buttons = append(buttons, newButton)
			do(message.Sender.ID, db.GetRecent(message.Sender.ID), text)
		} else if strings.HasPrefix(button[0].Text, message.Data) && message.Data != "Done" {
			buttons = append(buttons, newButton)
			undo(message.Sender.ID, db.GetRecent(message.Sender.ID), text)
		} else if message.Data != "Done" || strings.HasSuffix(button[0].Text, "✅") {
			buttons = append(buttons, []tb.InlineButton{button[0]})
		}
	}
	newKeyboard := tb.ReplyMarkup{InlineKeyboard: buttons}
	bot.EditReplyMarkup(message.Message, &newKeyboard)
}

func CreateMediaByThing(thing *constants.Thing, caption func(*constants.Thing) string) *tb.Photo {
	return &tb.Photo{
		File:    tb.File{FileID: thing.Photo},
		Caption: caption(thing),
		ParseMode: constants.ParseMode,
	}
}
