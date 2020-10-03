package states

import (
	"bot/pkg/constants"
	"bot/pkg/db"
	"github.com/sonichka1311/tgbotapi"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func Remind(bot *tgbotapi.BotAPI, id int64) {
	go bot.Send(tgbotapi.NewMessage(id, constants.Remind))
}

func Hello(bot *tgbotapi.BotAPI, id int64) {
	go bot.Send(tgbotapi.NewMessage(id, constants.Hello))
}

func Help(bot *tgbotapi.BotAPI, id int64) {
	go bot.Send(tgbotapi.NewMessage(id, constants.Help))
}

func SmthWrong(bot *tgbotapi.BotAPI, id int64) {
	go bot.Send(tgbotapi.NewMessage(id, constants.CommandNotFound))
}

func GetThing(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	id, _ := strconv.Atoi(message.Text[7:])
	thing := db.GetThing(message.From.ID, id)
	photo := tgbotapi.NewPhotoShare(message.Chat.ID, thing.Photo)
	photo.ParseMode = constants.ParseMode
	photo.Caption = constants.ThingText(thing)
	bot.Send(photo)
}

func Wardrobe(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	rows := db.GetAll(message.From.ID)
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
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.EmptyWardrobe))
	}
	txts := constants.SplitBigMsg(text)
	for _, txt := range txts {
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		msg.ParseMode = constants.ParseMode
		bot.Send(msg)
	}
}

func MakeDirty(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	id, _ := strconv.Atoi(message.Text[7:])
	db.MakeDirty(message.From.ID, id)
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.MarkThingTo(db.GetName(message.From.ID, id), false)+message.Text[7:])
	msg.ParseMode = constants.ParseMode
	bot.Send(msg)
}

func MakeClean(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	id, _ := strconv.Atoi(message.Text[7:])
	db.MakeClean(message.From.ID, id)
	msg := tgbotapi.NewMessage(message.Chat.ID, constants.MarkThingTo(db.GetName(message.From.ID, id), true)+message.Text[7:])
	msg.ParseMode = constants.ParseMode
	bot.Send(msg)
}

func Dirty(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	rows := db.GetDirty(message.From.ID)
	texts := make([]string, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe from database: %s\n", dbError.Error())
		}
		texts = append(texts, constants.ThingsText(&thing))
	}
	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Ура! Грязных вещей нет."))
	}
	txts := constants.SplitBigMsg(text)
	for _, txt := range txts {
		msg := tgbotapi.NewMessage(message.Chat.ID, txt)
		msg.ParseMode = constants.ParseMode
		bot.Send(msg)
	}
}

func GetByType(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message) {
	rows := db.GetByType(message.From.ID, message.Text[6:])
	text := ""
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Purity, &thing.Photo)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe by type from database: %s\n", dbError.Error())
		}
		text += "*" + thing.Name + "* (" + thing.Purity + ")\nПосмотреть: /thing\\_" + strconv.Itoa(thing.Id) + "\n"
		if thing.Purity == "dirty" {
			text += "Отметить чистым: /clean\\_" + strconv.Itoa(thing.Id) + "\n\n"
		} else {
			text += "Отметить грязным: /dirty\\_" + strconv.Itoa(thing.Id) + "\n\n"
		}
	}
	if len(text) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "В этой категории пока ничего нет. Добавь новую вещь: /upload"))
	}
	for len(text) > 4096 {
		ind := strings.LastIndex(text[:4096], "\n\n")
		msg := tgbotapi.NewMessage(message.Chat.ID, text[:ind])
		msg.ParseMode = "Markdown"
		bot.Send(msg)
		text = text[ind+2:]
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
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

func SendSomething(bot *tgbotapi.BotAPI, message *tgbotapi.Message, things []*constants.Thing, change bool, thingType string) {
	if len(things) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, constants.NoCleanThing[thingType]))
		return
	}
	thing := things[rand.Intn(len(things))]
	//log.Println(strings.Split(message.Caption, "\n")[0])
	if change {
		if len(things) > 1 {
			//log.Println(strings.Split(message.Caption, "\n")[0][len(thingType):])
			//log.Println(strings.Split(constants.Caption(thing), "\n")[0][len(thingType) + 2:])
			for strings.Split(message.Caption, "\n")[0][len(thingType):] == strings.Split(constants.Caption(thing), "\n")[0][len(thingType)+2:] {
				log.Println("The same thing occurred.")
				//log.Println(strings.Split(message.Caption, "\n")[0])
				//log.Println(strings.Split(constants.Caption(thing), "\n")[0])
				thing = things[rand.Intn(len(things))]
			}
		} else {
			msg := tgbotapi.NewEditMessageCaption(
				message.Chat.ID,
				message.MessageID,
				constants.Caption(thing)+"В твоем гардеробе есть только одна чистая вещь этого типа.",
			)
			msg.ParseMode = constants.ParseMode
			bot.Send(msg)
			return
		}
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		constants.ChangeButtons(
			thing.Type,
			constants.ChangeFrom[thing.Type],
			constants.ChangeTo[thing.Type],
		)...,
	)
	photo := tgbotapi.InputMediaPhoto{
		Media:     thing.Photo,
		ParseMode: constants.ParseMode,
		Caption:   constants.Caption(thing),
	}
	messageId := message.MessageID
	if change {
		msg := tgbotapi.NewEditMessagePhoto(message.Chat.ID, message.MessageID, photo)
		msg.ReplyMarkup = &keyboard
		log.Println(msg.Media.Caption)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewPhotoShare(message.Chat.ID, photo.Media)
		msg.ParseMode = photo.ParseMode
		msg.Caption = photo.Caption
		msg.ReplyMarkup = &keyboard
		log.Println(msg.Caption)
		m, _ := bot.Send(msg)
		messageId = m.MessageID
	}
	go func() {
		time.Sleep(time.Minute * 5)
		msg := tgbotapi.NewEditMessageCaption(
			message.Chat.ID,
			messageId,
			constants.Caption(thing)+"Время на изменение выбора истекло.",
		)
		msg.ParseMode = constants.ParseMode
		bot.Send(msg)
	}()
}

func ChangeThing(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message, types string, change bool) {
	topColor := db.GetTopColor(message.From.ID)
	bottomColor := db.GetBottomColor(message.From.ID)
	colors := map[string]string{
		"top":    topColor,
		"bottom": bottomColor,
	}
	season := db.GetSeason(message.From.ID)
	color := "any"
	if c, ok := colors[types]; ok {
		color = c
	}
	SendSomething(bot, message, GetByParams(db, message.From.ID, color, types, season), change, types)
}

func GetRandomThing(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message, types string) {
	rows := db.GetByType(message.From.ID, types)
	ans := make([]*constants.Thing, 0)
	for rows.Next() {
		var thing constants.Thing
		dbError := rows.Scan(&thing.Id, &thing.Name, &thing.Purity, &thing.Photo)
		if dbError != nil {
			log.Printf("Error while selecting all wardrobe by type from database: %s\n", dbError.Error())
		}
		if thing.Purity == "clean" {
			ans = append(ans, &thing)
		}
	}
	if len(ans) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "В этой категории пока ничего нет. Добавь новую вещь: /upload"))
		return
	}
	num := rand.Intn(len(ans))
	msg := tgbotapi.NewPhotoShare(message.Chat.ID, ans[num].Photo)
	msg.ParseMode = "Markdown"
	msg.Caption = "*" + ans[num].Name + "*\n" +
		"Посмотреть подробнее: /thing\\_" + strconv.Itoa(ans[num].Id) + "\n"
	log.Println(msg.Caption)
	bot.Send(msg)
}

func DeleteThing(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.Message, idStr string) {
	id, _ := strconv.Atoi(idStr)
	db.DeleteThing(message.From.ID, id)
	bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Вещь удалена. /help"))
}

func MultiCallback(bot *tgbotapi.BotAPI, db *db.Database, message *tgbotapi.CallbackQuery, do, undo func(id int, recent int, s string)) {
	constants.MutexMap[message.From.ID].Lock()
	defer constants.MutexMap[message.From.ID].Unlock()
	keyboard := message.Message.ReplyMarkup
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, button := range keyboard.InlineKeyboard {
		newButton := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(message.Data, button[0].Text)}
		text := strings.ToLower(strings.Split(button[0].Text, " ")[0])
		if strings.HasPrefix(message.Data, button[0].Text) && message.Data != "Done" {
			buttons = append(buttons, newButton)
			do(message.From.ID, db.GetRecent(message.From.ID), text)
		} else if strings.HasPrefix(button[0].Text, message.Data) && message.Data != "Done" {
			buttons = append(buttons, newButton)
			undo(message.From.ID, db.GetRecent(message.From.ID), text)
		} else if message.Data != "Done" || strings.HasSuffix(button[0].Text, "✅") {
			buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button[0]})
		}
	}
	msg := tgbotapi.NewEditMessageReplyMarkup(message.Message.Chat.ID, message.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup(buttons...))
	bot.Send(msg)
}
