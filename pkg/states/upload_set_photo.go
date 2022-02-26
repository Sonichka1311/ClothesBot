package states

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
	"bot/pkg/s3"
)

type UploadSetPhotoState struct {
	BaseState
}

func NewUploadSetPhotoState(bot *tb.Bot, db *db.Database, s3 *s3.S3) State {
	return &UploadSetPhotoState{BaseState: NewBase(bot, db, s3)}
}

func (s UploadSetPhotoState) Do(message *tb.Message) string {
	recent := s.db.GetUser(message.Sender.ID).LastFileID

	s.db.AddThing(message.Sender.ID, recent)
	go func() { s.uploadToS3(message.Photo.FileID, s.removeBackground(s.getRawData(message.Photo))) }()
	s.db.SetPhoto(message.Sender.ID, recent, message.Photo.FileID)

	s.bot.Send(message.Sender, constants.SendMeName)
	return UploadSetNameState{}.GetName()
}

func (s UploadSetPhotoState) GetName() string {
	return "uploadSetPhoto"
}

func (s UploadSetPhotoState) getRawData(photo *tb.Photo) []byte {
	file, _ := s.bot.GetFile(photo.MediaFile())
	defer file.Close()

	data, _ := io.ReadAll(file)
	return data
}

func (s UploadSetPhotoState) removeBackground(data []byte) []byte {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fileWriter, err := w.CreateFormFile("content", "file.jpeg")
	if err != nil {
		panic(err)
	}
	fileWriter.Write(data)
	w.Close()

	r, _ := http.NewRequest(
		http.MethodPost,
		"https://pixcut.wondershare.com/openapi/api/v1/matting/removebg",
		&b,
	)
	r.Header.Set("Content-Type", w.FormDataContentType())
	r.Header.Set("appkey", "")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println("Error when removing background:", err)
	}

	defer resp.Body.Close()
	res, _ := io.ReadAll(resp.Body)
	return res
}

func (s UploadSetPhotoState) uploadToS3(key string, data []byte) {
	err := s.s3.PutObject(key, data)
	if err != nil {
		log.Println("Error while upload file to S3:", err.Error())
	}
}
