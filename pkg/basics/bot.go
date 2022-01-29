package basics

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/states"
)

type Bot struct {
	BotToken string
	Bot      *tb.Bot
	DB       *db.Database
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

		settings := tb.Settings{
			Token:     b.BotToken,
			Poller:    &tb.LongPoller{Timeout: time.Second},
			ParseMode: constants.ParseMode,
		}

		bot, err := tb.NewBot(settings)
		if err != nil {
			log.Fatalf("NewBot failed: %s", err)
		}

		b.Bot = bot
		b.DB = &db.Database{DB: dbSql}

		bot.Handle(tb.OnCallback, b.HandleCallback)
		bot.Handle(tb.OnText, b.HandleMessage)
		bot.Handle(tb.OnPhoto, b.HandleMessage)
		bot.Start()
	}
	return nil
}

func (b *Bot) HandleMessage(message *tb.Message) {
	log.Println("Got message " + message.Text)

	if message.Text == "/start" {
		b.DB.AddUser(message.Sender.ID)
		states.Hello(b.Bot, message.Sender)
		return
	}

	if message.Text == "/end" {
		b.DB.UpdateState(message.Sender.ID, states.MainState{}.GetName())
		states.Remind(b.Bot, message.Sender)
		return
	}

	if message.Text == "/help" {
		b.DB.UpdateState(message.Sender.ID, states.MainState{}.GetName())
		states.Help(b.Bot, message.Sender)
		return
	}

	state := states.States[b.DB.GetState(message.Sender.ID)]
	b.DB.UpdateState(message.Sender.ID, state.Do(b.Bot, b.DB, message))
}

func (b *Bot) HandleCallback(message *tb.Callback) {
	log.Println("Got message " + message.Data)

	state := states.States[b.DB.GetState(message.Sender.ID)]

	if message.Data == "Done" {
		states.MultiCallback(b.Bot, b.DB, message, nil, nil)
		msg := message.Message
		msg.Sender = message.Sender
		b.DB.UpdateState(message.Sender.ID, state.Do(b.Bot, b.DB, msg))
		return
	}

	switch state {
	case states.UploadSetSeasonState{}:
		states.MultiCallback(b.Bot, b.DB, message, b.DB.SetSeason, b.DB.UnsetSeason)
	case states.MainState{}:
		splitMessage := strings.Split(message.Data, "_")

		message.Message.Sender = message.Sender
		if splitMessage[0] == "type" {
			states.ChangeThing(b.Bot, b.DB, message.Data, message.Message, splitMessage[1], false)
		} else {
			states.ChangeThing(b.Bot, b.DB, message.Data, message.Message, splitMessage[0], true)
		}
	default:
		button := constants.NewButton(message.Data, message.Data)
		keyboard := constants.NewKeyboard(button)
		b.Bot.EditReplyMarkup(message.Message, keyboard)

		b.DB.UpdateState(
			message.Sender.ID,
			state.Do(
				b.Bot, b.DB,
				&tb.Message{
					Text:   message.Data,
					Sender: &tb.User{ID: message.Sender.ID},
					Chat:   &tb.Chat{ID: message.Message.Chat.ID},
				},
			),
		)
	}
}
