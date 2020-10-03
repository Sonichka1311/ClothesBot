package basics

import (
	"bot/pkg/db"
	"bot/pkg/states"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sonichka1311/tgbotapi"
	"log"
	"strings"
	"time"
)

type Bot struct {
	BotToken string
	//WebURL 	string
	Bot *tgbotapi.BotAPI
	DB  *db.Database
}

func (b *Bot) StartBot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("Context is done.")
	default:
		dbSql, dbError := sql.Open("mysql", "root:guest@tcp(mysql:3306)/clothes?charset=utf8&interpolateParams=true")
		if dbError != nil {
			log.Fatalf("Cannot open database: %s", dbError.Error())
		}

		count := 0
		for err := dbSql.Ping(); err != nil; err = dbSql.Ping() {
			log.Printf("DB is unavailable: %s. Retry in 5 sec...", err)
			time.Sleep(time.Second * 5)
			if count > 10 {
				log.Fatal("DB is unavailable. ")
			}
			count++
		}

		bot, err := tgbotapi.NewBotAPI(b.BotToken)
		if err != nil {
			log.Fatalf("NewBotAPI failed: %s", err)
		}
		b.Bot = bot

		b.DB = &db.Database{DB: dbSql}

		var updates = tgbotapi.NewUpdate(0)
		updates.Timeout = 60
		ch, err := b.Bot.GetUpdatesChan(updates)
		for {
			select {
			case message := <-ch:
				go b.Do(&message)
			}
		}
	}
	return nil
}

func (b *Bot) Do(message *tgbotapi.Update) {
	if message.Message != nil {
		b.HandleMessage(message.Message)
	} else if message.CallbackQuery != nil {
		b.HandleCallback(message.CallbackQuery)
	}
}

func (b *Bot) HandleMessage(message *tgbotapi.Message) {
	log.Println("Got message " + message.Text)
	if message.Text == "/start" {
		b.DB.AddUser(message.From.ID)
		states.Hello(b.Bot, message.Chat.ID)
		return
	}
	if message.Text == "/end" {
		b.DB.UpdateState(message.From.ID, states.MainState{}.GetName())
		states.Remind(b.Bot, message.Chat.ID)
		return
	}
	if message.Text == "/help" {
		b.DB.UpdateState(message.From.ID, states.MainState{}.GetName())
		states.Help(b.Bot, message.Chat.ID)
		return
	}
	state := states.States[b.DB.GetState(message.From.ID)]
	b.DB.UpdateState(message.From.ID, state.Do(b.Bot, b.DB, message))
}

func (b *Bot) HandleCallback(message *tgbotapi.CallbackQuery) {
	state := states.States[b.DB.GetState(message.From.ID)]
	if message.Data == "Done" {
		states.MultiCallback(b.Bot, b.DB, message, nil, nil)
		b.DB.UpdateState(message.From.ID, state.Do(b.Bot, b.DB, message.Message))
		return
	}
	switch state {
	case states.SeasonState{}:
		states.MultiCallback(b.Bot, b.DB, message, b.DB.SetSeason, b.DB.UnsetSeason)
	case states.MainState{}:
		splitMessage := strings.Split(message.Data, "_")
		message.Message.From = message.From
		if splitMessage[0] == "type" {
			if splitMessage[1] == "comb" {
				states.ChangeThing(b.Bot, b.DB, message.Message, "combo", false)
			} else if splitMessage[1] == "sep" {
				states.ChangeThing(b.Bot, b.DB, message.Message, "top", false)
				states.ChangeThing(b.Bot, b.DB, message.Message, "bottom", false)
			}
		} else {
			states.ChangeThing(b.Bot, b.DB, message.Message, splitMessage[0], true)
		}
	default:
		button := tgbotapi.NewInlineKeyboardButtonData(message.Data, message.Data)
		keyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{button})
		msg := tgbotapi.NewEditMessageReplyMarkup(message.Message.Chat.ID, message.Message.MessageID, keyboard)
		b.Bot.Send(msg)

		b.DB.UpdateState(
			message.From.ID,
			state.Do(
				b.Bot, b.DB,
				&tgbotapi.Message{
					Text: message.Data,
					From: &tgbotapi.User{ID: message.From.ID},
					Chat: &tgbotapi.Chat{ID: message.Message.Chat.ID},
				},
			),
		)
	}
}
