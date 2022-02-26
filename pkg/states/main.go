package states

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
)

type MainState struct {
	BaseState
}

func NewMainState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &MainState{BaseState: NewBase(bot, db, s3)}
}

func (s MainState) Do(message *tb.Message) string {
	splitMessage := strings.Split(message.Text, "_")
	switch splitMessage[0] {
	case "/upload":
		go s.bot.Send(message.Sender, constants.SendMePhoto)
		return UploadSetPhotoState{}.GetName()
	case "/wardrobe":
		s.wardrobe(message)
		return s.GetName()
	case "/thing":
		s.getThing(message)
		return s.GetName()
	case "/look":
		go s.bot.Send(message.Sender, constants.WhatColor(strings.ToLower(constants.Top)), constants.ColorButtons(true))
		return LookSetTopColorState{}.GetName()
	case "/dirty":
		if len(splitMessage) > 1 {
			s.changePurity(message, false)
		} else {
			s.listDirty(message)
		}
		return s.GetName()
	case "/clean":
		s.changePurity(message, true)
		return s.GetName()
	case "/show":
		s.getByType(message)
		return s.GetName()
	case "/get":
		s.getRandomThing(message, splitMessage[1])
		return s.GetName()
	case "/delete":
		s.deleteThing(message, splitMessage[1], true)
		s.bot.Send(message.Sender, constants.Deleted)
		return s.GetName()
	//case "/all":
	//tb.NewMediaGroup()
	//return s.GetName()
	default:
		s.smthWrong(message.Sender)
		return s.GetName()
	}
}

func (s MainState) GetName() string {
	return "main"
}

func (s MainState) getThing(message *tb.Message) {
	id, _ := strconv.Atoi(message.Text[7:])
	thing := s.db.GetThing(message.Sender.ID, id)
	photo := s.createMediaByThing(thing, thing.Caption)
	s.bot.Send(message.Sender, photo)
}

func (s MainState) wardrobe(message *tb.Message) {
	things := s.db.GetThingsByUser(message.Sender.ID)
	texts := make([]string, 0)

	for _, thing := range things {
		texts = append(texts, thing.ListCaption())
	}

	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		s.bot.Send(message.Sender, constants.EmptyArray["wardrobe"])
	}

	s.sendBigMessage(message.Sender, text)
}

func (s MainState) changePurity(message *tb.Message, toClean bool) {
	id, _ := strconv.Atoi(message.Text[7:])

	if toClean {
		s.db.MakeClean(message.Sender.ID, id)
	} else {
		s.db.MakeDirty(message.Sender.ID, id)
	}

	s.bot.Send(message.Sender, constants.MarkThingTo(s.db.GetThing(message.Sender.ID, id).Name, toClean)+message.Text[7:])
}

func (s MainState) listDirty(message *tb.Message) {
	things := s.db.ListDirty(message.Sender.ID)
	texts := make([]string, 0)

	for _, thing := range things {
		texts = append(texts, thing.ListCaption())
	}

	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		s.bot.Send(message.Sender, constants.EmptyArray["dirty"])
	}

	s.sendBigMessage(message.Sender, text)
}

func (s MainState) getByType(message *tb.Message) {
	things := s.db.ListByType(message.Sender.ID, message.Text[6:])
	texts := make([]string, 0)

	for _, thing := range things {
		texts = append(texts, thing.ListCaption())
	}

	text := strings.Join(texts, "\n")
	if len(text) == 0 {
		s.bot.Send(message.Sender, constants.EmptyArray["by_type"])
	}

	s.sendBigMessage(message.Sender, text)
}


func (s MainState) getRandomThing(message *tb.Message, types string) {
	things := s.db.ListByType(message.Sender.ID, types)
	ans := make([]*constants.Thing, 0)

	for _, thing := range things {
		if thing.Purity == "clean" {
			ans = append(ans, thing)
		}
	}

	if len(ans) == 0 {
		s.bot.Send(message.Sender, constants.EmptyArray["random"])
		return
	}

	thing := ans[rand.Intn(len(ans))]
	photo := s.createMediaByThing(thing, thing.ShortCaption)
	s.bot.Send(message.Sender, photo)
	log.Println(photo.Caption)
}
