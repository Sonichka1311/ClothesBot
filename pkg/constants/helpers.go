package constants

import tb "gopkg.in/tucnak/telebot.v2"

func NewButton(text, data string) *tb.InlineButton {
	return &tb.InlineButton{
		Text:            text,
		Data:            data,
		InlineQueryChat: "",
	}
}

func NewKeyboard(buttons ...*tb.InlineButton) *tb.ReplyMarkup {
	keyboard := make([][]tb.InlineButton, 0)
	for _, button := range buttons {
		keyboard = append(keyboard, []tb.InlineButton{*button})
	}
	return &tb.ReplyMarkup{InlineKeyboard: keyboard}
}
