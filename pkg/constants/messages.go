package constants

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/utils"
)

const (
	Remind = "Нужна помощь? /help"
	Hello  = "Привет! Я могу помочь тебе выбрать образ на сегодня. Чтобы узнать, что я умею, нажми /help"
	Help   = "Вот, что я умею:\n" +
		"• /upload - добавить новую вещь\n" +
		"• /wardrobe - показать все вещи\n" +
		"• /show\\_top - показать весь верх\n" +
		"• /show\\_bottom - показать весь низ\n" +
		"• /show\\_combo - показать все комбинезоны и платья\n" +
		"• /show\\_outer - показать всю верхнюю одежду\n" +
		"• /show\\_shoes - показать всю обувь\n" +
		"• /look - сгенерировать образ\n" +
		"• /dirty - показать все грязные вещи\n" +
		"• /get\\_top - показать случайный верх\n" +
		"• /get\\_bottom - показать случайный низ\n" +
		"• /get\\_combo - показать случайный(ое) комбинезон или платье\n" +
		"• /get\\_outer - показать случайную верхнюю одежду\n" +
		"• /get\\_shoes - показать случайную обувь\n"
	CommandNotFound    = "Не нашел такую команду. Нужна помощь? /help"
	SendMePhoto        = "Пришли мне фотографию новой вещи. Если хочешь прервать добавление вещи, нажми /end"
	SendMeName         = "Круто! Как называется?"
	SendMeType         = "А что за вещь? Выбери нужный вариант:"
	SendMeSeason       = "Для какой она погоды? Выбери все подходящие варианты:"
	SendMeCombo        = "Эту вещь можно носить отдельно или в комбинации с другими? Выбери все подходящие варианты:"
	SendMeColor        = "Это тёмная или светлая вещь? Выбери нужный вариант:"
	WhatSeason         = "Какая сейчас погода? Выбери нужный вариант:"
	Cold               = "Холодно (ниже 5°)"
	Normal             = "Нормально (5° - 15°)"
	Warm               = "Тепло (15° - 23°)"
	Hot                = "Жарко (23°+)"
	Done               = "Готово"
	Top                = "Верх"
	Bottom             = "Низ"
	Combo              = "Комбинация"
	Outer              = "Верхняя одежда"
	Shoes              = "Обувь"
	Light              = "Светлое"
	Dark               = "Тёмное"
	Any                = "Любое"
	Sep                = "Раздельный"
	Comb               = "Совместный"
	Dirty              = "Грязное"
	Clean              = "Чистое"
	NeedSomethingClean = "Чтобы сгенерировать образ, нужен хотя бы один верх+низ или комбо в гардеробе."
	ParseMode          = "Markdown"
	EmptyWardrobe      = "Гардероб пока пуст :( Добавь новую вещь: /upload"
	MaxLength          = 4096
	Deleted            = "Вещь была удалена. /help"
	JustOneThing       = "В гардеробе только одна вещь такого типа.\n"
	TimeIsUp           = "Время для смены образа истекло.\n"
	UploadCancelled    = "Загрузка отменена."
)

var (
	Added = func(name string, id int) string {
		return fmt.Sprintf(
			`Вещь *\"%s\"* добавлена в гардероб.\n
			Показать полную информацию: /thing\\_%d\n
			Отметить как грязное: /dirty\\_%d\n
			Добавить новую вещь: /upload
		`, name, id, id)
	}

	NoThing = func(typeThing string) string {
		return fmt.Sprintf("В гардеробе нет подходящей чистой вещи типа %s.", typeThing)
	}

	WhatColor = func(typeThing string) string {
		return fmt.Sprintf("Хочешь %s конкретного цвета? Выбери подходящий вариант: ", typeThing)
	}

	SeasonButtons = func(needEnd bool) *tb.ReplyMarkup {
		cold := NewButton(Cold, Cold+" ✅")
		normal := NewButton(Normal, Normal+" ✅")
		warm := NewButton(Warm, Warm+" ✅")
		hot := NewButton(Hot, Hot+" ✅")
		end := NewButton(Done, Done)
		if needEnd {
			return NewKeyboard(cold, normal, warm, hot, end)
		}
		return NewKeyboard(cold, normal, warm, hot)
	}

	TypeButtons = func() *tb.ReplyMarkup {
		top := NewButton(Top, Top+" ✅")
		bottom := NewButton(Bottom, Bottom+" ✅")
		combo := NewButton(Combo, Combo+" ✅")
		outer := NewButton(Outer, Outer+" ✅")
		shoes := NewButton(Shoes, Shoes+" ✅")
		return NewKeyboard(top, bottom, combo, outer, shoes)
	}

	ColorButtons = func(needAny bool) *tb.ReplyMarkup {
		light := NewButton(Light, Light+" ✅")
		dark := NewButton(Dark, Dark+" ✅")
		any := NewButton(Any, Any+" ✅")
		if needAny {
			return NewKeyboard(any, light, dark)
		}
		return NewKeyboard(light, dark)
	}

	//ChangeButtons = func(this, fromType, toType string) *tb.ReplyMarkup {
	//	changeThis := NewButton("Change "+this, this)
	//	changeType := NewButton("Change "+fromType+" to "+toType, "type_"+toType)
	//	if len(fromType) == 0 || len(toType) == 0 {
	//		return NewKeyboard(changeThis)
	//	} else {
	//		return NewKeyboard(changeThis, changeType)
	//	}
	//}

	ChangeButton = func(this string) *tb.InlineButton {
		return NewButton("Change "+this, this)
	}

	ChangeTypeButton = func(fromType, toType string) *tb.InlineButton {
		return NewButton("Change "+fromType+" to "+toType, "type_"+toType)
	}

	MarkThingTo = func(name string, isClean bool) string {
		marked := Dirty
		offered := Clean
		if isClean {
			marked = Clean
			offered = Dirty
		}
		marked = strings.ToLower(marked)
		offered = strings.ToLower(offered)
		return fmt.Sprintf("Вещь *%s* Отмечена как %s. Отметить как %s: /%s\\_", name, marked, offered, utils.ToEng(offered))
	}

	NoCleanThing = map[string]string{ // type to text
		"top":    NoThing("top"),
		"bottom": NoThing("bottom"),
		"combo":  NoThing("combo or dress"),
		"outer":  NoThing("outerwear"),
		"shoes":  NoThing("shoes"),
		"sep":    NeedSomethingClean,
	}

	ChangeFrom = map[string]string{
		"top":    Sep,
		"bottom": Sep,
		"combo":  Comb,
		"outer":  "",
		"shoes":  "",
	}

	ChangeTo = map[string]string{
		"top":    Comb,
		"bottom": Comb,
		"combo":  Sep,
		"outer":  "",
		"shoes":  "",
	}

	EmptyArray = map[string]string{ // func name to text
		"dirty":    "Ура! Нет грязных вещей.",
		"by_type":  "Нет ничего в этой категории :( Добавь новую вещь: /upload",
		"wardrobe": EmptyWardrobe,
		"random":   "В этой категории пока ничего нет. Добавь новую вещь: /upload",
	}
)
