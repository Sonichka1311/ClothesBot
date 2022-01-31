package states

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/aws/aws-sdk-go/service/s3"
	tb "gopkg.in/tucnak/telebot.v2"

	"bot/pkg/constants"
	"bot/pkg/db"
)

type UploadSetPhotoState struct{}

func (s UploadSetPhotoState) Do(bot *tb.Bot, db *db.Database, s3Client *s3.S3, message *tb.Message) string {
	recent := db.GetUser(message.Sender.ID).LastFileID

	db.AddThing(message.Sender.ID, recent)
	s.uploadToS3(s3Client, message.Photo.FileID, s.removeBackground(s.getRawData(bot, message.Photo)))
	db.SetPhoto(message.Sender.ID, recent, message.Photo.FileID)

	bot.Send(message.Sender, constants.SendMeName)
	return UploadSetNameState{}.GetName()
}

func (s UploadSetPhotoState) GetName() string {
	return "uploadSetPhoto"
}

func (s UploadSetPhotoState) getRawData(bot *tb.Bot, photo *tb.Photo) []byte {
	file, _ := bot.GetFile(photo.MediaFile())
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
	r.Header.Set("appkey", "<>")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Println("Error when removing background:", err)
	}

	defer resp.Body.Close()
	res, _ := io.ReadAll(resp.Body)
	return res
}

func (s UploadSetPhotoState) uploadToS3(
	s3Client *s3.S3,
	key string,
	data []byte,
) {
	bucketName := "<>"
	body := bytes.NewReader(data)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   body,
	})
	if err != nil {
		log.Println("Error while upload file to S3:", err.Error())
	}
}
