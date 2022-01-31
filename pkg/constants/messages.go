package constants

import (
	"fmt"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	Remind = "Do you need /help?"
	Hello  = "Hey! I can help you choose your today's outfit! To find out what I can do, click /help"
	Help   = "That's what I can do:\n" +
		"• /upload - add new thing\n" +
		"• /wardrobe - show all things\n" +
		"• /show\\_top - show all tops\n" +
		"• /show\\_bottom - show all bottoms\n" +
		"• /show\\_combo - show all combos and dresses\n" +
		"• /show\\_outer - show all outerwear\n" +
		"• /show\\_shoes - show all shoes\n" +
		"• /look - generate outfit\n" +
		"• /dirty - show all dirty things\n" +
		"• /get\\_top - generate one top\n" +
		"• /get\\_bottom - generate one bottom\n" +
		"• /get\\_combo - generate one combo or dress\n" +
		"• /get\\_outer - generate one outerwear\n" +
		"• /get\\_shoes - generate one shoes\n"
	CommandNotFound    = "I didn't find that command. Do you need /help?"
	SendMePhoto        = "Send me photo of new thing. If you want to stop adding new thing, click /end"
	SendMeName         = "How do you call it?"
	SendMeType         = "What type is it of? Tap on the right variant:"
	SendMeSeason       = "What weather is it for? Choose all appropriate variants:"
	SendMeColor        = "Is it light or dark thing? Tap on the right variant:"
	WhatSeason         = "What weather is it now? Tap on the right variant: "
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
	NeedSomethingClean = "You need to have at least once clean top and bottom or combo to generate outfit"
	ParseMode          = "Markdown"
	EmptyWardrobe      = "Your wardrobe is empty. Add new thing: /upload"
	MaxLength          = 4096
	Deleted 		   = "Thing was deleted. /help"
	JustOneThing       = "There are just one clean thing of this type.\n"
	TimeIsUp           = "The time for change outfit is up.\n"
)

var (
	Added = func(name string, id int) string {
		return fmt.Sprintf("Thing *\"%s\"* was added to your wardrobe.\n" +
			"See more: /thing\\_%d\nMark dirty: /dirty\\_%d\nAdd new thing: /upload", name, id, id)
	}

	NoThing = func(typeThing string) string {
		return fmt.Sprintf("You don't have appropriate clean %s.", typeThing)
	}

	WhatColor = func(typeThing string) string {
		return fmt.Sprintf("Do you want %s of specific color? Tap on the right variant: ", typeThing)
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
		return NewButton("Change " + this, this)
	}

	ChangeTypeButton = func(fromType, toType string) *tb.InlineButton {
		return NewButton("Change "+fromType+" to "+toType, "type_"+toType)
	}

	Caption = func(thing *Thing) string {
		return fmt.Sprintf(
			"*%s*: %s\nShow more: /thing\\_%d\n",
			strings.ToUpper(thing.Type[0:1])+thing.Type[1:], thing.Name, thing.ID,
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
			thing.Color, thing.Purity, toCond, toCond, thing.ID, thing.ID,
		)
	}

	MarkThingTo = func(name string, isClean bool) string {
		marked := Dirty
		offered := Clean
		if isClean {
			marked = Clean
			offered = Dirty
		}
		return fmt.Sprintf("Thing *%s* marked as %s. Mark as %s: /%s\\_", name, marked, offered, offered)
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

	EmptyArray = map[string]string { // func name to text
		"dirty": "Hurray! There are no dirty things.",
		"by_type": "There are nothing in this category. Add new thing: /upload",
		"wardrobe": EmptyWardrobe,
		"random": "В этой категории пока ничего нет. Добавь новую вещь: /upload",
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
