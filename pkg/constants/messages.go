package constants

import (
	"fmt"
	"github.com/sonichka1311/tgbotapi"
	"strings"
)

const (
	Remind = "Напомнить, что я умею? /help"
	Hello  = "Привет! Я твой помощник в выборе одежды на сегодня! Чтобы узнать, что я умею, нажми /help"
	Help   = "Вот, что я умею:\n" +
		"* /upload - добавить новую вещь\n" +
		"* /wardrobe - узнать, что есть в твоем гардеробе\n" +
		"* /show_top - посмотреть весь верх твоего гардероба\n" +
		"* /show_bottom - посмотреть весь низ твоего гардероба\n" +
		"* /show_combo - посмотреть все комбинезоны и платья\n" +
		"* /show_outer - посмотреть верхнюю одежду\n" +
		"* /show_shoes - посмотреть обувь\n" +
		"* /look - подобрать одежду на сегодня\n" +
		"* /dirty - посмотреть всю грязную одежду (отметить одежду грязной/чистой можно по командам из показа одежды)\n" +
		"* /get_top - Подобрать случайный верх\n" +
		"* /get_bottom - подобрать случайный низ\n" +
		"* /get_combo - подобрать случайный комбинезон или платье\n" +
		"* /get_outer - подобрать случайную верхнюю одежду\n" +
		"* /get_shoes - подобрать случайную обувь"
	CommandNotFound    = "Ты ввел что-то не так. Напомнить, что я умею? /help"
	SendMePhoto        = "Отправь мне фотографию вещи, которую хочешь добавить. Если хочешь прервать добавление вещи, нажми /end"
	SendMeName         = "Как называется эта вещь? Напиши желаемое отображение этой вещи в свободной форме"
	SendMeType         = "Что это за вещь? Нажми на подходящий вариант:"
	SendMeSeason       = "Для какого сезона эта вещь? Выбери все подходящие варианты:"
	SendMeColor        = "Это светлая или темная вещь? Нажми на подходящий вариант:"
	WhatTopColor       = "Возьмем верх определенного цвета? Нажми на подходящий вариант: "
	WhatBottomColor    = "Возьмем низ определенного цвета? Нажми на подходящий вариант: "
	WhatSeason         = "Для какого сезона смотрим наряд? Выбери подходящий вариант: "
	Cold               = "Cold (under 5°)"
	Normal             = "Normal (5° - 15°)"
	Warm               = "Warm (15° - 23°)"
	Hot                = "Hot (23°+)"
	Done               = "Done"
	Top                = "Top"
	Bottom             = "Bottom"
	Combo              = "Combo"
	Outer              = "Outer"
	Shoes              = "Shoes"
	Light              = "Light"
	Dark               = "Dark"
	Any                = "Any"
	Sep                = "sep"
	Comb               = "comb"
	Dirty              = "dirty"
	Clean              = "clean"
	NeedSomethingClean = "Чтобы собрать лук, нужнен хотя бы один чистый верх и один чистый низ или одно чистое комбо"
	NoTop              = "В твоем гардеробе нет подходящего чистого верха"
	NoBottom           = "В твоем гардеробе нет подходящего чистого низа"
	NoCombo            = "В твоем гардеробе нет подходящих чистых комбо"
	NoShoes            = "В твоем гардеробе нет подходящей чистой обуви"
	NoOuter            = "В твоем гардеробе нет подходящей чистой верхней одежды"
	ParseMode          = "Markdown"
	EmptyWardrobe      = "Пока в твоем гардеробе ничего нет. Добавь новую вещь: /upload"
	MaxLength          = 4096
)

var (
	Added = func(name string, id int) string {
		return fmt.Sprintf("Вещь \"%s\" добавлена в гардероб.\nПосмотреть /thing_%d\nОтметить грязной /dirty_%d\nДобавить новую вещь: /upload", name, id, id)
	}

	SeasonButtons = func(needEnd bool) [][]tgbotapi.InlineKeyboardButton {
		cold := tgbotapi.NewInlineKeyboardButtonData(Cold, Cold+" ✅")
		normal := tgbotapi.NewInlineKeyboardButtonData(Normal, Normal+" ✅")
		warm := tgbotapi.NewInlineKeyboardButtonData(Warm, Warm+" ✅")
		hot := tgbotapi.NewInlineKeyboardButtonData(Hot, Hot+" ✅")
		end := tgbotapi.NewInlineKeyboardButtonData(Done, Done)
		if needEnd {
			return [][]tgbotapi.InlineKeyboardButton{{cold}, {normal}, {warm}, {hot}, {end}}
		}
		return [][]tgbotapi.InlineKeyboardButton{{cold}, {normal}, {warm}, {hot}}
	}

	TypeButtons = func() [][]tgbotapi.InlineKeyboardButton {
		top := tgbotapi.NewInlineKeyboardButtonData(Top, Top+" ✅")
		bottom := tgbotapi.NewInlineKeyboardButtonData(Bottom, Bottom+" ✅")
		combo := tgbotapi.NewInlineKeyboardButtonData(Combo, Combo+" ✅")
		outer := tgbotapi.NewInlineKeyboardButtonData(Outer, Outer+" ✅")
		shoes := tgbotapi.NewInlineKeyboardButtonData(Shoes, Shoes+" ✅")
		return [][]tgbotapi.InlineKeyboardButton{{top}, {bottom}, {combo}, {outer}, {shoes}}
	}

	ColorButtons = func(needAny bool) [][]tgbotapi.InlineKeyboardButton {
		light := tgbotapi.NewInlineKeyboardButtonData(Light, Light+" ✅")
		dark := tgbotapi.NewInlineKeyboardButtonData(Dark, Dark+" ✅")
		any := tgbotapi.NewInlineKeyboardButtonData(Any, Any+" ✅")
		if needAny {
			return [][]tgbotapi.InlineKeyboardButton{{any}, {light}, {dark}}
		}
		return [][]tgbotapi.InlineKeyboardButton{{light}, {dark}}
	}

	ChangeButtons = func(this, fromType, toType string) [][]tgbotapi.InlineKeyboardButton {
		changeThis := tgbotapi.NewInlineKeyboardButtonData("Change "+this, this)
		changeType := tgbotapi.NewInlineKeyboardButtonData("Change "+fromType+" to "+toType, "type_"+toType)
		if len(fromType) == 0 || len(toType) == 0 {
			return [][]tgbotapi.InlineKeyboardButton{{changeThis}}
		} else {
			return [][]tgbotapi.InlineKeyboardButton{{changeThis}, {changeType}}
		}
	}

	Caption = func(thing *Thing) string {
		return fmt.Sprintf(
			"*%s*: %s\nShow more: /thing\\_%d\n",
			strings.ToUpper(thing.Type[0:1])+thing.Type[1:], thing.Name, thing.Id,
		)
	}

	ThingText = func(thing *Thing) string {
		seasons := make([]string, 0)
		if thing.Cold {
			seasons = append(seasons, "cold")
		}
		if thing.Normal {
			seasons = append(seasons, "normal")
		}
		if thing.Warm {
			seasons = append(seasons, "warm")
		}
		if thing.Hot {
			seasons = append(seasons, "hot")
		}
		var toCond string
		if thing.Purity == "dirty" {
			toCond = Clean
		} else {
			toCond = Dirty
		}

		return fmt.Sprintf(
			"*%s*\n"+
				"*Type:* %s\n"+
				"*Season:* %v\n"+
				"*Color:* %s\n"+
				"*Condition:* %s\n"+
				"Mark thing *%s*: /%s\\_%d\n"+
				"Delete thing: /delete\\_%d\n",
			thing.Name, thing.Type, strings.Join(seasons, ", "),
			thing.Color, thing.Purity, toCond, toCond, thing.Id, thing.Id,
		)
	}

	ThingsText = func(thing *Thing) string {
		var toCond string
		if thing.Purity == "dirty" {
			toCond = Clean
		} else {
			toCond = Dirty
		}
		return fmt.Sprintf(
			"*%s* (%s)\n"+
				"See more: /thing\\_%d\n"+
				"Mark thing *%s*: /%s\\_%d\n",
			thing.Name, thing.Purity, thing.Id, toCond, toCond, thing.Id,
		)
	}

	MarkThingTo = func(name string, isClean bool) string {
		marked := Dirty
		offered := Clean
		if isClean {
			marked = Clean
			offered = Dirty
		}
		return fmt.Sprintf("Thing *%s* marked as %s. Mark as %s: /%s_", name, marked, offered, offered)
	}

	NoCleanThing = map[string]string{ // type to text
		"top":    NoTop,
		"bottom": NoBottom,
		"combo":  NoCombo,
		"outer":  NoOuter,
		"shoes":  NoShoes,
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

	SplitBigMsg = func(text string) []string {
		texts := make([]string, 0)
		for len(text) > 0 {
			if len(text) < MaxLength {
				texts = append(texts, text)
				break
			}
			ind := strings.LastIndex(text[:MaxLength], "\n\n")
			texts = append(texts, text[:ind])
			text = text[ind+2:]
		}
		return texts
	}
)
