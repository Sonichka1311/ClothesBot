package basics

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/states"
)

type Bot struct {
	BotToken string
	Bot      *tb.Bot
	DB       *db.Database
	S3       *s3.S3
}

func (b *Bot) StartBot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("Context is done.")
	default:
		dbSql, dbError := sql.Open(
			"postgres",
			`
				host=<>
				port=6432
				dbname=clothes
				user=<>
				password=<>
				sslmode=verify-full
				sslrootcert=/Users/<>/.postgresql/root.crt
			`,
		)
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

		gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: dbSql}), &gorm.Config{})
		if err != nil {
			log.Println("can't open gorm db", err)
			return err
		}

		b.DB = &db.Database{DB: dbSql, Gorm: gormDB}

		s3Endpoint := "https://storage.yandexcloud.net"
		s3Session, _ := session.NewSession(&aws.Config{
			Region:   aws.String("us-west-2"),
			Endpoint: &s3Endpoint,
		})

		b.S3 = s3.New(s3Session)

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
		b.DB.CreateUser(message.Sender.ID)
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

	state := states.States[b.DB.GetUser(message.Sender.ID).State]
	b.DB.UpdateState(message.Sender.ID, state.Do(b.Bot, b.DB, b.S3, message))
}

func (b *Bot) HandleCallback(message *tb.Callback) {
	log.Println("Got message " + message.Data)

	state := states.States[b.DB.GetUser(message.Sender.ID).State]

	if message.Data == "Done" {
		states.MultiCallback(b.Bot, b.DB, message, nil, nil)
		msg := message.Message
		msg.Sender = message.Sender
		b.DB.UpdateState(message.Sender.ID, state.Do(b.Bot, b.DB, b.S3, msg))
		return
	}

	switch state {
	case states.UploadSetSeasonState{}:
		states.MultiCallback(b.Bot, b.DB, message, b.DB.SetSeason, b.DB.UnsetSeason)
	case states.MainState{}:
		splitMessage := strings.Split(message.Data, "_")

		message.Message.Sender = message.Sender
		if splitMessage[0] == "type" {
			states.ChangeThing(b.Bot, b.DB, b.S3, message.Data, message.Message, splitMessage[1], false)
		} else {
			states.ChangeThing(b.Bot, b.DB, b.S3, message.Data, message.Message, splitMessage[0], true)
		}
	default:
		button := constants.NewButton(message.Data, message.Data)
		keyboard := constants.NewKeyboard(button)
		b.Bot.EditReplyMarkup(message.Message, keyboard)

		b.DB.UpdateState(
			message.Sender.ID,
			state.Do(
				b.Bot, b.DB, b.S3,
				&tb.Message{
					Text:   message.Data,
					Sender: &tb.User{ID: message.Sender.ID},
					Chat:   &tb.Chat{ID: message.Message.Chat.ID},
				},
			),
		)
	}
}
