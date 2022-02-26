package basics

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
	"bot/pkg/states"
)

type Bot struct {
	BotToken string
	Bot      *tb.Bot
	DB       *db.Database
	S3       *s3.S3
	States   *states.StateFabric
}

func (b *Bot) StartBot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Println("Context is done.")
	default:
		dbSql, dbError := sql.Open(
			"postgres",
			`
				host=
				port=6432
				dbname=
				user=
				password=
				sslmode=verify-full
				sslrootcert=
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
		b.S3 = s3.NewS3()

		b.States = states.NewStateFabric(b.Bot, b.DB, b.S3)

		bot.Handle(tb.OnCallback, b.HandleCallback)
		bot.Handle(tb.OnText, b.HandleMessage)
		bot.Handle(tb.OnPhoto, b.HandleMessage)
		bot.Start()
	}
	return nil
}

func (b *Bot) HandleMessage(message *tb.Message) {
	log.Println("Got message " + message.Text)

	baseState := b.States.NewState("base")
	nextState := baseState.Do(message)
	if nextState != baseState.GetName() {
		b.DB.UpdateState(message.Sender.ID, nextState)
		return
	}

	state := b.States.NewState(b.DB.GetUser(message.Sender.ID).State)
	b.DB.UpdateState(message.Sender.ID, state.Do(message))
}

func (b *Bot) HandleCallback(message *tb.Callback) {
	log.Println("Got message " + message.Data)

	baseState := states.NewBase(b.Bot, b.DB, b.S3)
	state := b.States.NewState(b.DB.GetUser(message.Sender.ID).State)

	if message.Data == constants.Done {
		baseState.MultiCallback(message, nil, nil)
		msg := message.Message
		msg.Sender = message.Sender
		b.DB.UpdateState(message.Sender.ID, state.Do(msg))
		return
	}

	switch state.GetName() {
	case states.UploadSetSeasonState{}.GetName():
		baseState.MultiCallback(message, b.DB.SetSeason, b.DB.UnsetSeason)
	case states.MainState{}.GetName():
		splitMessage := strings.Split(message.Data, "_")

		message.Message.Sender = message.Sender
		if splitMessage[0] == "type" {
			baseState.ChangeThing(message.Data, message.Message, splitMessage[1], false)
		} else {
			baseState.ChangeThing(message.Data, message.Message, splitMessage[0], true)
		}
	default:
		button := constants.NewButton(message.Data, message.Data)
		keyboard := constants.NewKeyboard(button)
		b.Bot.EditReplyMarkup(message.Message, keyboard)

		b.DB.UpdateState(
			message.Sender.ID,
			state.Do(
				&tb.Message{
					Text:   message.Data,
					Sender: &tb.User{ID: message.Sender.ID},
					Chat:   &tb.Chat{ID: message.Message.Chat.ID},
				},
			),
		)
	}
}
